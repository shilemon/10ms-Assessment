## ops-2

Name: Emon Shil
Email: emonshil;.htuc@gmail.com
Time taken: around 6–7 hours (deployment, debugging, monitoring setup, RCA documentation)
Environment used (kind/minikube): kind (Kubernetes in Docker)
Any assumptions:

Monitoring stack (Prometheus + Grafana + metrics-server) is required for validating HPA, pod health, and resource usage.

Docker images must be rebuilt and pushed to a registry for KIND to pull correctly.

Port configuration and probe endpoints are assumed to be customizable as part of the assessment.

Grafana dashboards and metrics verification are acceptable evidence of monitoring setup.
