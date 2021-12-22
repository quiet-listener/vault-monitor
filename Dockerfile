FROM golang:1.16 as builder

WORKDIR /vault_monitor
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
ENV GOPROXY=direct
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY vault/ vault/


# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o vault_monitor main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static
WORKDIR /
COPY --from=builder /vault_monitor/vault_monitor .

ENTRYPOINT ["/vault_monitor"]
