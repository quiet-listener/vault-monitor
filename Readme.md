# vault-monitor
Exposes prometheus endpoint with vault specific metrics from given cluster.Currently generated metrics:

`vault_ping{cluster="mow-dev-eu-central-1",account="manowar_dev",instance="10.84.0.137:8080",job="kubernetes-pods",kubernetes_namespace="default",kubernetes_pod_name="vault-monitor",name="vault-monitor"} 1`

| Values | Interpretation |
|--|--|
| 0 | Unavailable (not able to ping vault from cluster. Check logs of pod for possible reasons) |
| 1 | Available (able to ping vault from cluster) |

| Tags available | Values example|
|--|--|
|cluster|"mow-dev-eu-central-1"|
|account|"manowar_dev"|
|instance|"10.84.0.137:8080"|
|job|"kubernetes-pods"|
|kubernetes_namespace|"default"|
|kubernetes_pod_name|"vault-monitor"|
|name|"vault-monitor"|
