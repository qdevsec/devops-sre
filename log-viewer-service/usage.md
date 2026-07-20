## logviewer

service that allows you to view log sources you wire up on a web page instead of being restricted to terminal

- high availability
- sidecars you can configure [syslog, auditd]
- REST endpoints
- components for IaC deployment

## spinning up

##### create image
`$ cd log-viewer-service/containers/logviewer`

`$ npm init -y`

`$ npm install express`

`$ eval $(minikube docker-env)`

`$ docker build -t logviewer:latest .`

##### apply ymls

`$ kubectl apply -f k8s/config-logviewer-html.yaml `

`$ kubectl apply -f k8s/deployment-ha-logviewer.yaml `

`$ kubectl apply -f k8s/service-ha-logviewer.yaml`

`$ kubectl apply -f k8s/prometheus-config.yaml`

`$ kubectl apply -f k8s/prometheus.yaml`

`$ kubectl apply -f k8s/grafana.yaml`

##### view UIs

`$ minikube service ha-logviewer-svc -n web-ha`

`$ minikube service prometheus-svc -n web-ha`

`$ minikube service grafana-svc -n web-ha`

grafana login

`login admin/admin`

- test/monitoring has started dashboards to import


