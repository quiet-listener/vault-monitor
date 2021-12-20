package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	vm "vault_monitor/vault"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	csi_cluster = os.Getenv("CSI_CLUSTER")
	csi_env     = os.Getenv("CSI_ENV")
	gauge       = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "vault",
			Subsystem:   "core",
			Name:        "ping",
			Help:        "Liveness Check of Vault from Given Cluster",
			ConstLabels: prometheus.Labels{"account": csi_env, "cluster": csi_cluster},
		})
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "prom_request_time",
		Help:        "Time it has taken to retrieve the metrics",
		ConstLabels: prometheus.Labels{"env": csi_env, "cluster": csi_cluster},
	}, []string{"time"})
	router.HandleFunc("/", homeLink)
	router.Handle("/metrics", newHandlerWithHistogram(promhttp.Handler(), histogramVec))
	prometheus.Register(histogramVec)
	prometheus.MustRegister(gauge)
	go func() {
		for {
			live := float64(1)
			_, err := vm.VaultPing()
			if err != nil {
				log.Print(err)
				live = float64(0)
			} else {
				live = float64(1)
			}
			gauge.Set(live)
			time.Sleep(time.Minute)
		}
	}()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func newHandlerWithHistogram(handler http.Handler, histogram *prometheus.HistogramVec) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		status := http.StatusOK

		defer func() {
			histogram.WithLabelValues(fmt.Sprintf("%d", status)).Observe(time.Since(start).Seconds())
		}()

		if req.Method == http.MethodGet {
			handler.ServeHTTP(w, req)
			return
		}
		status = http.StatusBadRequest

		w.WriteHeader(status)
	})
}
