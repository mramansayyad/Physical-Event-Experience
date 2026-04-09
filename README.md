# Physical Event Experience Platform 🏟️

### **Architectural Mesh for High-Concurrency Stadium Intelligence**
**Engineering Lead:** Aman Sayyad, Founder & CEO, Aman Tech Innovations  
**Version:** `v1.1.0-SECURED` (Hardened)  
**Tech Stack:** Go 1.25.9, Google Cloud Run, Redis, Firestore, Docker (Distroless)

---

## 🏗️ 10x Engineering Architecture
This platform implements a **Hexagonal Serverless Mesh** (Ports & Adapters) to ensure complete decoupling of business logic from infrastructure.

* **Core Domain:** Isolated logic for crowd heatmap calculation and real-time gate rerouting.
* **Adapters:** High-performance implementations for GCP Firestore and Redis with connection pooling.
* **Transport:** Hardened HTTP layer with "Always-Fail" validation logic.

## 🛡️ Security & Zero-Trust Profile
* **Hardened Runtime:** Deployed using **Google Distroless (Non-Root)** images to minimize the attack surface.
* **Vulnerability Scanning:** Automated **SecOps Sweep** via `govulncheck` integrated into the CI/CD pipeline.
* **Secret Management:** Zero use of `.env` files in production; all credentials are dynamically fetched from **GCP Secret Manager**.
* **Patch Management:** Fully patched against critical standard library vulnerabilities (GO-2026-4870, 4865, 4947, 4946) by enforcing **Go 1.25.9**.

## ⚡ Performance & Efficiency
* **Circuit Breakers:** Implemented `sony/gobreaker` patterns to prevent cascading failures during database latency spikes.
* **Connection Pooling:** Optimized Redis buffers with a pool size of 100 to handle massive telemetry surges.
* **Worker Pools:** Concurrency-safe telemetry processing using buffered channels to prevent OOM (Out of Memory) events.

## 📊 Observability & Reliability
* **Probing:** Native `/healthz` (Liveness) and `/readyz` (Readiness) endpoints for automated cloud self-healing.
* **Graceful Shutdown:** Interception of `SIGTERM` signals to ensure zero data loss for inflight telemetry during deployments.
* **Structured Logging:** JSON-based logging mapped to GCP Cloud Logging severity levels.
* **Tracing:** OpenTelemetry integration for sub-second request tracing across the hexagonal mesh.

## 🚀 Deployment
Deployed on **Google Cloud Run** using a multi-stage Docker build for millisecond cold starts and infinite auto-scaling.

```bash
# Production Deployment Sequence
gcloud run deploy stadium-backend --source . --region us-central1 --allow-unauthenticated
