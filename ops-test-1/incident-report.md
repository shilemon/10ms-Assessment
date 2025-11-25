 Incident Report


## Executive Summary

The SRE application experienced a complete service outage due to multiple configuration misalignments across the application, Docker, and Kubernetes layers. All pods were in a non-ready state, preventing any traffic from reaching the application. Additionally, the application had performance issues causing 10-second response delays.

Duration: Continuous until fixed  
Impact: 100% service unavailability  
Users Affected: All users  


## Impact

### User Impact
- Availability: 0% - Service completely unavailable
- Users Affected: All users attempting to access the service
- **Duration:** Continuous until configuration fixes applied

### Business Impact
- Complete service outage
- No successful health checks
- Zero request processing capability
- Potential reputation damage

### Technical Impact
- All pod replicas in NotReady state
- Service endpoints empty
- Health check endpoints failing
- Application performance degraded (when it would have been accessible)

---

## Root Cause Analysis

### Primary Root Cause

**Configuration Misalignment Across Multiple Layers**

The incident was caused by a cascade of configuration errors:

1. **Application Layer Issues**
   - Health check endpoint returning HTTP 500 (incorrect status)
   - Artificial 10-second delay in main endpoint (`time.sleep(10)`)
   - Application correctly configured to run on port 8080

2. **Kubernetes Configuration Issues**
   - Container port defined as 80 (should be 8080)
   - Readiness probe checking port 80 (should be 8080)
   - Missing imagePullPolicy for `latest` tag
   - No resource limits defined
   - No liveness probe configured

3. **Docker Configuration Issues**
   - Wrong port exposed (80 instead of 8080)
   - Inefficient image build (no layer caching)
   - Running as root user (security risk)
   - Using full Python image instead of slim variant

### Why It Happened

1. **Port Mismatch:** The application runs on port 8080, but the Kubernetes deployment was configured with containerPort 80, and the readiness probe was checking port 80. This created a situation where:
   - The application was listening on 8080
   - Kubernetes was expecting it on 80
   - Health checks couldn't connect (connection refused)

2. **Health Endpoint Failure:** Even when the correct port was used, the health endpoint returned HTTP 500 instead of 200, causing all readiness checks to fail.

3. **Cascading Failure:** Because readiness probes never passed:
   - Pods never entered "Ready" state
   - Service never added pod IPs to endpoints
   - No traffic could reach the application
   - Service appeared completely down

### Contributing Factors

- Lack of local testing before deployment
- No validation of health endpoint status codes
- Missing integration tests for port configurations
- No automated checks for configuration consistency

---

## What Was Fixed

### 1. Application Code (`app/main.py`)

**Before:**
```python
@app.route('/')
def home():
    time.sleep(10)  # ❌ Artificial delay
    return "Hello from SRE App!", 200

@app.route('/healthz')
def health():
    return "Unhealthy", 500  #  Wrong status code
```

**After:**
```python
from flask import Flask, jsonify
import time

app = Flask(__name__)

@app.route('/')
def home():
    return jsonify({"message": "Hello from SRE Assessment!", "status": "healthy"})

@app.route('/healthz')
def health():
    return jsonify({"status": "healthy"}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)de
```

### 2. Dockerfile

**Improvements:**
- Changed to `python:3.9-slim` for smaller image
- Implemented proper layer caching (requirements first)
- Added non-root user for security
- Corrected exposed port to 8080
- Added HEALTHCHECK directive

### 3. Kubernetes Deployment (`k8s/deployment.yaml`)

**Before:**
```yaml
containers:
- name: sre-app
  image: sre-candidate:latest
  ports:
  - containerPort: 80  # ❌ Wrong port
  readinessProbe:
    httpGet:
      path: /healthz
      port: 80  # ❌ Wrong port
```

**After:**
```yaml
containers:
- name: sre-app
  image: emon110852/sre-assessment:v1.0
  imagePullPolicy: Always  # ✅ Added
  ports:
  - containerPort: 8080  # ✅ Correct port
  readinessProbe:
    httpGet:
      path: /healthz
      port: 8080  # ✅ Correct port
    initialDelaySeconds: 5
    periodSeconds: 5
  livenessProbe:  # ✅ Added
    httpGet:
      path: /healthz
      port: 8080
    initialDelaySeconds: 10
    periodSeconds: 10
  resources:  # ✅ Added
    requests:
      memory: "64Mi"
      cpu: "100m"
    limits:
      memory: "128Mi"
      cpu: "200m"
```

### 4. Service Configuration

**Verified correct configuration:**
```yaml
ports:
- protocol: TCP
  port: 80          # External port
  targetPort: 8080  # ✅ Correctly points to container port
```

---

## Detection and Response

### How It Was Detected
- Pods continuously showing NotReady status
- `kubectl get pods` showing 0/1 ready state
- Events showing readiness probe failures
- Service had no endpoints

### Investigation Process
1. Checked pod status: `kubectl get pods`
2. Examined pod details: `kubectl describe pod <pod-name>`
3. Reviewed pod logs: `kubectl logs <pod-name>`
4. Analyzed events: `kubectl get events --sort-by='.lastTimestamp'`
5. Inspected service endpoints: `kubectl get endpoints`
6. Tested application locally with Docker
7. Identified port mismatches and health check issues

### Resolution Steps
1. Fixed application code (removed delay, fixed status code)
2. Rebuilt Docker image with corrections
3. Updated Kubernetes manifests with correct ports
4. Added missing configurations (resources, liveness probe)
5. Pushed new image to registry
6. Redeployed application
7. Verified all pods reached Ready state
8. Tested service endpoints

---

## Preventive Measures

### Immediate Actions (Completed)

1. ✅ **Configuration Validation**
   - Aligned all port configurations across layers
   - Verified health check endpoints return 200
   - Added explicit imagePullPolicy

2. ✅ **Enhanced Monitoring**
   - Added liveness probe for container health
   - Configured resource limits to prevent exhaustion
   - Improved probe timing parameters

3. ✅ **Security Improvements**
   - Running as non-root user
   - Added security context
   - Using read-only root filesystem where possible

### Short-Term Actions (Recommended)

1. **Pre-Deployment Testing**
   - Implement local testing procedure before deployment
   - Create integration tests for health endpoints
   - Test Docker images locally before pushing

2. **Automation**
   - Add Makefile for consistent build/deploy process
   - Implement GitHub Actions CI for automated testing
   - Add linting for Kubernetes YAML files

3. **Documentation**
   - Document port requirements
   - Create deployment runbook
   - Document health check requirements


---


### Error Messages

```
Warning  Unhealthy  readiness probe failed: HTTP probe failed with statuscode: 500

Warning  Unhealthy  readiness probe failed: 
Get "http://3.0.148.238/:80/healthz": dial tcp 3.0.148.238:80: 
connect: connection refused

Warning  BackOff  Back-off restarting failed container
```

### Verification Commands

```bash
# Verify pods are ready
kubectl get pods -l app=sre
# Expected: All pods showing 1/1 READY

# Check service endpoints
kubectl get endpoints sre-service
# Expected: Multiple pod IPs listed

# Test health endpoint
kubectl port-forward service/sre-service 8080:80
curl http://127.0.0.1:8080/healthz
# Expected: "Healthy" with HTTP 200

```


- DevOps team (deployment configuration)
- SRE team (monitoring and alerting)
- Management (process improvements)

**Action Items Review Date:** [Set date for 2 weeks out]
