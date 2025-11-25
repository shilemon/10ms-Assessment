# Troubleshooting Documentation

# Application Fixes

### Issue 1: Slow Home Endpoint
Problem**: The `/` endpoint had a `time.sleep(5)` call causing 5-second delays.
Fix: Removed the unnecessary sleep statement.
Impact: Response time improved from 5+ seconds to <100ms.

### Issue 2: Health Check Status Code
Problem: `/healthz` endpoint was returning status code 500 instead of 200.
Fix: Changed return statement from `return jsonify(...), 500` to `return jsonify(...), 200`.
Impact: Readiness probes now pass successfully.

## Dockerfile Fixes

### Issue 1: Build Failure
Problem: Requirements not being installed, wrong base image.
Fix: 
- Used `python:3.9-slim` base image
- Added proper COPY and RUN commands for requirements.txt
 Impact: Image builds successfully.

### Issue 2: Inefficient Layering
Problem: Application code copied before requirements, breaking cache.
Fix: Copied requirements.txt first, then installed dependencies, then copied app code.
Impact: Faster rebuild times when only code changes.

### Issue 3: Security
Problem: Running as root user.
Fix: Created non-root user `appuser` and switched to it.
**Impact**: Improved container security posture.

## Kubernetes Fixes

### Issue 1: Image Pull Error
Problem: `imagePullPolicy: Always` with local image not in registry.
Fix: Changed to `imagePullPolicy: IfNotPresent` and loaded image into kind cluster.
**Commands**:
```bash
kind load docker-image emon110852/sre-assessment:latest --name sre-test
```

### Issue 2: Readiness Probe Failures
**Problem**: 
- Wrong probe path or port
- Health endpoint returning wrong status code
Fix: 
- Corrected probe path to `/healthz`
- Ensured port matches container port (8080)
- Fixed application to return 200 status
  Impact: Pods become ready and receive traffic.

### Issue 3: Service Not Routing
Problem: Selector labels in Service didn't match Deployment labels.
Fix: Ensured both use `app: sre-assessment` label.
Impact: Service properly routes traffic to pods.

### Issue 4: Port Mismatch
Problem: Service targetPort didn't match container port.
**Fix**: Set `targetPort: 8080` to match container's exposed port.
Impact: Traffic properly forwarded to application.




## Log Analysis

### Root Cause of Readiness Failures
1. Health endpoint returning 500 status code
2. Application response too slow (5+ seconds) exceeding probe timeout
3. Probe configuration issues (wrong path/port)

### Root Cause of Latency
- Unnecessary `time.sleep(5)` in home endpoint
- No performance optimization in application code

### Permanent Fixes Applied
1. Removed sleep delays from application code
2. Fixed health endpoint to return 200 status
3. Configured appropriate probe timeouts and initial delays
4. Added resource limits to prevent resource exhaustion

## Commands Used

### Docker Commands
```bash
docker build -t emon110852/sre-assessment:latest .
docker run -p 8080:8080 sre-assessment:latest
docker images | grep sre-assessment
```

### Kind Commands
```bash
kind create cluster --name sre-test
kind load docker-image emon110852/sre-assessment:latest --name sre-test
kind get clusters
```

### Kubernetes Commands
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl get pods -w
kubectl describe pod <pod-name>
kubectl logs -l app=sre-assessment
kubectl get svc
kubectl port-forward service/sre-assessment-service 8080:80
kubectl get endpoints
```

### Testing Commands
```bash
curl http://localhost:8080/
curl http://localhost:8080/healthz
```

## What I Learned

1. **Importance of Health Checks**: Properly configured health checks are crucial for zero-downtime deployments
2. **Docker Layer Caching**: Ordering Dockerfile commands properly significantly improves build times
3. **Label Selectors**: Kubernetes service routing depends on exact label matches
4. **Local Development**: Using kind/minikube requires special image loading steps
5. **Probe Configuration**: Timeouts and delays must account for application startup time
6. **Resource Management**: Setting requests/limits prevents resource contention

<img width="1371" height="242" alt="Screenshot 2025-11-25 215333" src="https://github.com/user-attachments/assets/82bdbc68-78b6-4ba9-81b2-a73cd5119b4a" />


## Create CI

- Automated validation on each change: The workflow compiles  Python, builds the Docker image, and smoke‑tests / and /healthz. This catches mistakes early (syntax errors, broken endpoints, mis‑built images).
- Repeatability: CI runs in a clean GitHub runner every time, ensuring  app builds and starts from scratch—not just on  local machine.
- Confidence for Kubernetes rollout: A container that passes basic health checks in CI is far more likely to pass readiness/liveness probes in cluster.
- Fast feedback for PRs: Contributors see immediate pass/fail on their pull requests before merging, reducing regressions.
- Optional publishing: When  uncomment the Docker Hub steps and add secrets, CI can publish emon110852/sre-assessment:latest automatically on pushes to main, making deployments consistent and quick.

