
## Spinning up

$ kubectl apply -f k8s/namespaces.yaml

1. Add OpenSearch Helm repo
$ helm repo add opensearch https://opensearch-project.github.io/helm-charts/
$ helm repo update

2. Deploy OpenSearch
$ helm install opensearch opensearch/opensearch -n security -f k8s/opensearch-values.yaml

3. Deploy OpenSearch Dashboards (using default values for now)
$ helm install opensearch-dashboards opensearch/opensearch-dashboards -n security

Watch pods until running

$ kubectl get pods -n security -w

when everything port-forward

$ kubectl port-forward svc/opensearch-dashboards 5601:5601 -n security

browse to 

http://localhost:5601

username: admin

login info is in opensearch-values.yaml


## Telemetry

Linux - auditd, journald, Fluent Bit

Windows - Sysmon, Windows Event Logs, Winlogbeat


Deploy fluent bit in k8s

$ helm repo add fluent https://fluent.github.io/helm-charts

$ helm repo update

$ helm install fluent-bit fluent/fluent-bit -n security --create-namespace

Host value in fluent-bit-values.yaml depends on service name created by opensearch helm chart

run the following to verify

$ kubectl get svc -n security

reinstall or upgrade fluent bit

$ helm upgrade fluent-bit fluent/fluent-bit \
    -n security \
    -f k8s/fluent-bit-values.yaml


May need to do this if on a linux host

$ sudo ufw allow 49900/tcp

$ minikube mount --port=49900 /var/log/audit:/var/log/audit


Verify data is arriving

$ kubectl logs -n security daemonset/fluent-bit

port forward opensearch

$ kubectl port-forward svc/opensearch-cluster-master 9200:9200 -n security

curl -k -u 'admin:<password>' https://localhost:9200/_cat/indices?v

curl -k -u 'admin:<YOUR_PASSWORD>' https://localhost:9200

Watch the logs

$ kubectl logs -n security daemonset/fluent-bit -f

## Detections

install sigma

can use to convert sigma detection rules

--
1st way to add detection to opensearch

$ sigma convert -t opensearch_lucene --without-pipeline kubernetes-opensearch-detection-siem/detections/<name of detection>

example result that is added to dashboard

(key:command_execution AND type:EXECVE) AND (a0:(*whoami* OR *uname*))

add detection to UI

http://localhost:5601 --> login --> dashboards --> Alerting --> Monitors --> Create Monitor 

Monitor Type: Per query monitor

Define Query: Extraction Query Editor

```
{
    "size": 0,
    "query": {
        "query_string": {
            "query" : "(key:command_execution AND type:EXECVE) AND (a0:(*whoami* OR *uname*))"
        }
    }
}
```

--
2nd way

with pysigma export to opensearch monitor json payload

sigma convert \
  -t opensearch_lucene \
  -f monitor_rule \
  --without-pipeline \
  kubernetes-opensearch-detection-siem/detections/T1082/linux_system_discovery.yaml \
  -o monitor.json

post to opensearch cluster

curl -XPOST "https://localhost:9200/_plugins/_alerting/monitors" \
     -H 'Content-Type: application/json' \
     -u "admin:your_password" \
     -d @monitor.json

## Redteaming

PowerShell on Ubuntu
$ sudo apt-get install -y wget apt-transport-https software-properties-common
$ wget -q "https://packages.microsoft.com/config/ubuntu/$(lsb_release -rs)/packages-microsoft-prod.deb"
$ sudo dpkg -i packages-microsoft-prod.deb
$ sudo apt-get update
$ sudo apt-get install -y powershell

add runner

$ pwsh

add YAML parser and the execution framework module

Install-Module -Name invoke-atomicredteam,powershell-yaml -Scope CurrentUser -Force

Download the execution framework and the Atomics test files

IEX (IWR 'https://raw.githubusercontent.com/redcanaryco/invoke-atomicredteam/master/install-atomicredteam.ps1' -UseBasicParsing); Install-AtomicRedTeam -getAtomics -Force

Import the module for  session

Import-Module ~/AtomicRedTeam/invoke-atomicredteam/Invoke-AtomicRedTeam.psd1 -Force

View details of a technique

Invoke-AtomicTest T1007 -ShowDetailsBrief

run test

Invoke-AtomicTest T1007

clean up

Invoke-AtomicTest T1007 -Cleanup

exit

#######
Notes

don't edit critical files like kube-apiserver.yaml, instead try workarounds like

minikube start \
  --extra-config=apiserver.audit-policy-file=/etc/kubernetes/audit/audit-policy.yaml \
  --extra-config=apiserver.audit-log-path=/var/log/kubernetes/audit.log


if you do decide to edit a key file make backups outside /etc/kubernetes/manifests/