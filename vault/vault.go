package vault

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	api "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"golang.org/x/net/http2"
)

func getK8sAuth() (*api.Client, error) {
	config := &api.Config{
		HttpClient: cleanhttp.DefaultPooledClient(),
	}
	config.HttpClient.Timeout = time.Second * 60

	transport := config.HttpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		config.Error = err
		return nil, fmt.Errorf("error while configuring tls %+v", err)
	}
	if err := config.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("error occured ReadEnvironment %+v", err)
	}

	config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	config.Backoff = retryablehttp.LinearJitterBackoff
	config.MaxRetries = 2

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}
	k8sAuth, err := auth.NewKubernetesAuth(
		os.Getenv("CSI_NAMESPACE")+"."+os.Getenv("CSI_SACC"),
		auth.WithServiceAccountTokenPath("/var/run/secrets/kubernetes.io/serviceaccount/token"),
		auth.WithMountPath(os.Getenv("CSI_CLUSTER")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Kubernetes auth method: %w", err)
	}
	authInfo, err := client.Auth().Login(context.TODO(), k8sAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}
	return client, nil
}

func VaultPing() (string, error) {
	client, err := getK8sAuth()
	if err != nil {

		return "", err
	}
	secret, err := client.Logical().Read("kv/data/public/test")
	if err != nil {
		return "", fmt.Errorf("unable to read secret: %w", err)
	}
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("data type assertion failed: %T %#v", secret.Data["data"], secret.Data["data"])
	}
	key := "answer"
	value, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("value type assertion failed: %T %#v", data[key], data[key])
	}
	return value, nil
}
