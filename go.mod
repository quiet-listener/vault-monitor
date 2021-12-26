//module github.infra.cloudera.com/thunderhead/datadog-iamauth-secret-creater

go 1.16

require (
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/vault/api v1.3.0
	github.com/hashicorp/vault/api/auth/kubernetes v0.1.0

	github.com/gorilla/mux v1.8.0
	github.com/prometheus/client_golang v1.11.0
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f
	github.com/hashicorp/go-retryablehttp v0.7.0
)

module vault_monitor
