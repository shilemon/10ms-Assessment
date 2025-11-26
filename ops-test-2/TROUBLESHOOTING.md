1. Dockerfile Fixes
1.1 Missing CA Certificates

The application could not access HTTPS endpoints because the container image lacked the system certificate bundle.
Installing the CA certificates resolved TLS and health-check failures.

1.2 Incorrect Exposed Port

The Dockerfile declared an incorrect port.
The application actually listens on port 8080, so the exposed port was corrected to match it.
This prevented Kubernetes readiness and liveness checks from failing.

1.3 Multi-Stage Build Issues

The original Dockerfile was not optimized.
It was rewritten into a proper multi-stage build structure, reducing the image size and ensuring that only the compiled application is included in the final layer.

1.4 Incorrect Copy Paths and Build Commands

Paths inside the Dockerfile were reorganized to ensure the application compiles and runs correctly.

2. Kubernetes Deployment Fixes
2.1 Incorrect Application Port in Deployment

The Deployment manifest was using the wrong application port.
It was corrected so Kubernetes would send traffic to the correct container port.

2.2 Readiness and Liveness Probes Failing

Health checks were referencing incorrect ports and incorrect paths.
These were corrected so the pod would stabilize and avoid crash loops.

2.3 Image Pull Policy

The pull policy was changed because KIND clusters cannot always fetch remote images quickly.
Updating the policy allowed Kubernetes to use local cached images more effectively.

2.4 Resource Requests and Limits Too Low

The application and sidecar had insufficient CPU and memory.
This caused throttling and HPA metrics instability.
Resource requests and limits were increased to stable and reasonable values.

3. Sidecar Fixes
3.1 Port Conflict Between App and Sidecar

Both the main app and the sidecar container were using the same port, which caused immediate pod failures.
The sidecar was assigned its own dedicated port.
Health checks were also adjusted to reference the new port.

3.2 Sidecar Health Endpoint

The sidecar container was missing a proper health endpoint reference, which prevented readiness and liveness from functioning.
The correct health path and port were configured.

4. Image Registry Fixes
4.1 Images Pushed Correctly to Docker Hub

Your images (sre-app2 and sidecar) were verified as publicly available.
Deployment references were updated to match the correct image names and tags.

4.2 KIND Local Image Syncing

Since KIND clusters require special handling to use local images, the deployment and pull policies were aligned accordingly.

5. Kind Cluster Issues
5.1 Wrong Cluster Context

You had multiple KIND clusters with similar names.
Your active cluster was identified as "sre-assess2".
The context was switched so all operations targeted the correct cluster.

5.2 Cluster Info Verification

Cluster connectivity and endpoints were validated to ensure manifests were applied to the correct environment.

6. Monitoring Setup Fixes
6.1 Prometheus

Prometheus was configured to collect metrics from:

the main application

the sidecar

Kubernetes metrics-server

6.2 Grafana

Dashboards were prepared for:

Pod health overview

HPA scaling behavior

Application latency and request rate

Sidecar metrics

6.3 HPA Metrics

Horizontal Pod Autoscaler was validated and scaled based on CPU or custom metrics.

6.4 Pod Health Visualization

Grafana dashboards were fully connected to Prometheus and displayed pod lifecycle, restarts, resource usage, and HTTP health checks.

✅ Summary of All Fixes Applied
Dockerfile

CA certificates added

Correct exposed port applied

Multi-stage build fixed

Build layers optimized

Paths and build steps corrected

Kubernetes

Correct container port

Probes fixed

Resource limits increased

Sidecar port conflict resolved

Image pull policy fixed

Cluster

Correct KIND cluster selected

Images pulled properly

Deployment connected to correct images

Monitoring

Prometheus installed

Grafana installed

HPA validated

Pod health dashboards working
