1. Summary

A service outage and monitoring failure occurred while deploying and observing workloads in the sre-assess2 cluster.
The incident was caused by a combination of:

misconfigured container ports

incorrect probes

a faulty Dockerfile

a sidecar port conflict

and incomplete monitoring configuration (Prometheus, Grafana, HPA metrics)

As a result:

application pods were stuck in CrashLoopBackOff

probes failed continuously

HPA metrics were missing

Grafana dashboards were blank

Prometheus scraped zero pod metrics

sidecar proxy was unreachable

All issues were investigated and resolved through multiple fixes at the container, Kubernetes, and monitoring levels.

2. Severity

Severity Level: SEV-2
Impact:

Application could not become Ready

Continuous restart loops

Sidecar failed due to port conflict

Readiness/Liveness probes failing

HPA unable to scale

Prometheus targets down

Grafana showing no pod/node metrics

Pod health panels stuck in No Data

3. Symptoms
Application-Level

Readiness and Liveness probes failing

Pods stuck in CrashLoopBackoff

“Development server” warnings inside production container

Wrong image pulled or not pulled due to imagePullPolicy

Sidecar-Level

Port conflict with main container

Health endpoint (/healthz) called repeatedly but unreachable

Monitoring-Level

Prometheus targets “DOWN”

Grafana dashboards missing CPU/Memory metrics

HPA showing unknown metrics

Pod health visualizations empty

What I observed in Grafana (screenshots added):

✔ Node metrics available
✔ Pod CPU/Memory metrics restored
✔ Cluster health panels working
✔ Prometheus scrape working

Included Screenshots:

Node resource usage

Pod health metrics

Prometheus target status

CPU utilization graphs

Cluster overview graphs

(These correspond to the images you attached.)

5. Root Cause Analysis (RCA)
Primary Root Cause

❌ Application containerPort did not match actual running port, causing
→ probes to fail
→ pod never becoming Ready
→ sidecar unable to connect

Secondary Causes

Sidecar binding to duplicate port

Dockerfile missing:

CA certificates

correct working directory

correct exposed port

Wrong probe path (/health vs /healthz)

Inaccurate resource limits causing throttling

KIND cluster using wrong imagePullPolicy

Prometheus scraping misconfigured

Metrics-server incompatible with cluster

Grafana datasource not applied correctly

6. Fixes Applied
Container Fixes

Rebuilt Dockerfile (proper multi-stage build)

Added CA certificates

Exposed correct application port

Cleaned build context

Kubernetes Fixes

Corrected containerPort: XXXX

Corrected readiness/liveness probe:

proper path

proper port

proper initialDelaySeconds

Updated resource limits

Resolved port conflict with sidecar

Set correct imagePullPolicy: IfNotPresent

Monitoring Fixes

Reinstalled kube-prometheus-stack

Fixed ServiceMonitor selectors

Added metrics-server with correct TLS flags

Repaired Grafana datasource

Restored Pod + Node dashboards

Validated HPA metrics

7. Impact Analysis

Application unavailable for ~1 hour

No data loss

HPA unable to scale → performance tests blocked

Dashboards unusable initially

Prometheus failed to collect metrics

Slow debugging due to misleading probe logs

Final Resolution

After fixing:

Container ports

Dockerfile

Sidecar ports

Probe configuration

Prometheus and Grafana configuration

metrics-server

HPA setup

The cluster became Stable, Healthy, and Fully Monitored.

<img width="1918" height="970" alt="image" src="https://github.com/user-attachments/assets/801c3c62-a9d6-45fe-95c6-1b67f44c009a" />

<img width="1910" height="963" alt="image" src="https://github.com/user-attachments/assets/1f61cf97-7dab-4c31-b902-733e623a325d" />

<img width="1916" height="870" alt="image" src="https://github.com/user-attachments/assets/aab609af-c2e7-4862-919c-827b798c77bc" />

<img width="1917" height="867" alt="image" src="https://github.com/user-attachments/assets/229fab3a-c357-4160-a61a-eb1a3900e85a" />

<img width="1919" height="880" alt="image" src="https://github.com/user-attachments/assets/71632d25-462c-4a13-8151-d3179c841fa0" />

<img width="1912" height="839" alt="image" src="https://github.com/user-attachments/assets/1c6eae82-1563-45c1-a274-3ce7742919eb" />




