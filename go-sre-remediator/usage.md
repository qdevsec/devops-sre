golang sre Project
- loads config.yaml that has rules rules (scripts to run based on alert)
- can build and run as container
- slackweb hook integration for posting formatted message to slack
- kubernetes deployment yaml
- observability (grafana)  

# start go server

$ cd /to/project/directory

$ go run main.go

Send fake json webhook payload

This is an alert that is sent, the system is set to take action when it sees this alert
It runs a script ./remediate_disk.sh 

The config rule specifies if DiskRunningFull alert then run the remediate_disk.sh
  - alert_name: "DiskRunningFull"
    script: "./remediate_disk.sh"

```
curl -X POST http://localhost:8080/webhook \
-H "Content-Type: application/json" \
-d '{
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "DiskRunningFull",
        "severity": "critical",
        "instance": "localhost"
      },
      "annotations": {
        "summary": "Disk space low on root partition"
      }
    }
  ]
}'
```

# Docker container

running container

$ cd /to/project/dir/where/dockerfile/is

$ docker build -t go-sre-remediator:v1 .

$ docker run -d -p 8080:8080 --name sre-engine go-sre-remediator:v1

with container running can run the following to see logs

$ docker logs -f sre-engine

After changes rebuild

# Slackweb hook

update const slackWebhookURL in main.go with url

uncomment slack integration code in the executescript() function

# Kubernetes deployment

use ConfigMap to inject config.yaml dynamically this allows for rules
to be updates in K8s without needing to rebuild the docker image

apply manifest to K8s cluster
$ kubectl apply -f deployment.yaml

verify pods are running
$ kubectl get pods -l app=auto-remediator

if you have the image local run

$ eval $(minikube docker-env)

then 

$ docker build -t go-sre-remediator:v1 .

to temporarily override local terminal docker command env vars so localhost docker daemon talks to docker daemon in minikube
with the docker build being ran again image is compiled inside minikube internal storage cache 

can start a rolling restart so kubernetes picks it up immediate
$ kubectl rollout restart deployment/sre-auto-remediator

# Grafana
install helm

$ helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
$ helm repo update

add stack to k8s cluster

helm install promo-stack prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false

view infrastructure

$ kubectl get pods -n monitoring --watch

apply the pod monitor

$ kubectl apply -f pod-monitor.yaml

$ kubectl port-forward pod/prometheus-promo-stack-kube-prometheu-prometheus-0 -n monitoring 9090:9090

view at
http://localhost:9090/targets

Port forward to Grafana (route Grafana's web UI to local port 3000)

$ kubectl port-forward pod/promo-stack-grafana-558bc95476-82cz8 -n monitoring 3000:3000

Get Admin Credentials

username: admin

$ kubectl get secret --namespace monitoring promo-stack-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo

view grafana at

http://localhost:3000