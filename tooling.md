# Comprehensive SRE & DevOps Tools Catalog

Site Reliability Engineering (SRE) and DevOps tools by category.

---

## 1. Infrastructure as Code (IaC) & Provisioning
*Tools used to define, spin up, and manage cloud infrastructure using declarative configurations.*

* **Terraform:** The industry standard for platform-agnostic, declarative infrastructure provisioning using HashiCorp Configuration Language (HCL).
* **OpenTofu:** The open-source, community-driven fork of Terraform, managed under the Linux Foundation.
* **Pulumi:** An Infrastructure as Code engine that allows engineers to write, deploy, and manage infrastructure using native programming languages (Python, TypeScript, Go, C#).
* **AWS CloudFormation / Azure Bicep / Google Cloud Deployment Manager:** Cloud-vendor native provisioning engines optimized for their respective cloud eco-systems.
* **Crossplane:** A Kubernetes-native control plane framework that enables the orchestration of infrastructure and external cloud services directly via Kubernetes Custom Resources (CRDs).

---

## 2. Configuration Management & GitOps
*Tools used to maintain the internal state of target machines and automate deployment synchronization via declarative Git repositories.*

* **Argo CD:** A highly popular, declarative GitOps continuous delivery tool designed specifically for Kubernetes applications.
* **Flux:** A CNCF secure GitOps delivery engine that natively reconciles Kubernetes clusters against Git/Helm repositories.
* **Ansible:** An agentless, SSH-based automation and configuration management engine for system provisioning and application deployments.
* **Chef / Puppet:** Configuration management platforms utilizing agent-based architectures for managing massive enterprise bare-metal or VM fleets.

---

## 3. Monitoring & Observability (The Telemetry Pillars)
*Tools built to collect, aggregate, and visualize metrics, logs, and distributed traces for deep architectural visibility.*

* **Prometheus & Grafana:** The ubiquitous open-source monitoring stack. Prometheus functions as a time-series pull-based metric engine, while Grafana acts as the visualization and dashboarding layer.
* **OpenTelemetry (OTel):** A vendor-neutral CNCF observability framework providing a standardized set of APIs, SDKs, and tools to generate, collect, and export telemetry data (Traces, Metrics, Logs).
* **Datadog / Dynatrace / New Relic:** Premium, all-in-one commercial SaaS observability platforms offering unified APM, logging, synthetics, and infrastructure monitoring.
* **Thanos / Cortex / Mimir:** Distributed, long-term metric storage engines designed to scale open-source Prometheus deployments across multiple clusters with global querying capabilities.
* **HyperDX / OpenObserve:** Modern, high-performance open-source log, trace, and metric aggregation alternatives built on unified columnar database engines (such as ClickHouse).

---

## 4. Incident Management, Alerting & On-Call
*Platforms used to manage schedules, coordinate incident response teams, automate alerting workflows, and minimize Mean Time to Resolution (MTTR).*

* **PagerDuty & Opsgenie:** Enterprise-ready alerting engines that manage complex on-call schedules, escalation pathways, and incident routing.
* **Rootly / incident.io:** Modern incident response hubs that automate the creation of Slack war-rooms, video conferencing bridge links, and structured post-mortem documents directly from chat interfaces.
* **Keep:** An open-core AIOps alert management platform designed to centralize, deduplicate, and correlate alerts from disparate monitoring tools prior to notifying on-call staff.
* **Alerta:** A lightweight, highly scalable centralized monitoring console for unified alert viewing. Developer friendly

---

## 5. Runbook Automation & Self-Healing
*Automated execution environments designed to reduce "toil" by triggering scripts, remediation jobs, and synthetic checks during incidents.*

* **PagerDuty Runbook Automation (formerly Rundeck):** A platform for codifying operational processes, turning standard operating procedures (SOPs) into self-service, push-button actions.
* **StackStorm:** Event-driven automation engine (IF-THIS-THEN-THAT for DevOps) that automatically triggers scripts, orchestrations, and healing workflows in response to monitoring alerts.
* **Upright:** An open-source synthetic monitoring and browser-automation tool designed to run continuous end-to-end user path tests natively from inside your Kubernetes cluster.

---

## 6. Continuous Integration & Continuous Delivery (CI/CD)
*Frameworks engineered to compile, test, package, and safely ship application binaries and services to runtime environments.*

* **GitHub Actions / GitLab CI/CD:** Modern, repository-native pipeline tools integrated seamlessly into the developers' Git workflow.
* **Harness:** An enterprise-focused continuous delivery platform emphasizing AI-driven automated canary rollouts, continuous verification, and automatic rollbacks.
* **Jenkins:** The classic, highly extensible, and community-backed self-hosted automation server with thousands of plugins.
* **Tekton:** A powerful, highly flexible Kubernetes-native CI/CD pipeline engine running workloads directly as containerized steps.

---

## 7. Chaos Engineering & Resilience Testing
*Tooling built to deliberately inject real-world faults (e.g., network latency, container death, disk stress) to verify systemic fault tolerance.*

* **Gremlin:** A managed SaaS platform built for safe, secure, and structured enterprise-scale chaos engineering experiments.
* **Chaos Mesh:** A highly visual, open-source cloud-native chaos engine built for Kubernetes, offering rich orchestration of network, stress, and file system faults.
* **LitmusChaos:** A community-led cloud-native framework focusing on a workflow-driven approach to testing resilience throughout the software development lifecycle.

---

## 8. Runtime Orchestration & Service Mesh
*The runtime layer executing distributed container applications and facilitating secure, observable communication between services.*

* **Kubernetes (K8s):** The universal operating system for orchestrating containerized application workloads at scale.
* **Istio / Linkerd:** Service meshes that manage mutual TLS (mTLS) encryption, secure traffic routing, canary splitting, and platform-level telemetry between microservices without modifying application code.
* **Nomad:** A highly lightweight, simple orchestrator by HashiCorp designed to schedule both containerized workloads and legacy non-containerized binaries.

---

## 9. DevSecOps & Compliance
*Tools built to scan application code, dependency libraries, container images, and host machines for vulnerabilities.*

* **Trivy / Aqua Security:** Rapid vulnerability scanners for container images, file systems, Git repositories, and IaC configurations.
* **Snyk:** Developer-first security tooling integrated directly into the local workspace and CI/CD pipelines to scan open-source dependencies, containers, and IaC code.
* **Lynis:** An open-source security auditing and hardening engine designed to inspect local Unix/Linux operating system configurations.