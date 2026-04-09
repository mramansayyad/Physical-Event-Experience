# Physical Event Experience Platform 🏟️

### **Architectural Mesh for High-Concurrency Stadium Intelligence**
**Engineering Lead:** Aman Sayyad  
**Organization:** Aman Tech Innovations  
**Version:** `v1.1.0-SECURED` (Hardened)  
**Tech Stack:** Go 1.25.9, Google Cloud Run, Redis, Firestore, Docker (Distroless)

---

## 🎯 Chosen Vertical: Physical Event Experience
This platform targets the **Smart Stadium and Large-Scale Event Management** vertical. It addresses the critical challenge of managing massive crowd density and real-time navigation during peak traffic periods, such as halftime rushes or emergency egress, within environments like Parul University’s venues.

---

## 🧠 Approach and Logic
The system follows a **10x Engineering Blueprint** centered on high-concurrency and "Zero-Trust" security.

### **Hexagonal Architecture (Ports & Adapters)**
We implemented a **Hexagonal Architecture** to strictly delineate core business logic from infrastructural dependencies.
* **Core Domain**: Isolated logic for crowd heatmap calculation and real-time gate rerouting algorithms.
* **Input Ports (Inbound)**: Decoupled HTTP transport layer with strict validation for incoming telemetry.
* **Output Ports (Outbound)**: Segregated interfaces for persistence (Firestore) and high-speed buffering (Redis).

### **Resiliency Patterns**
* **Circuit Breakers**: Implementation of `sony/gobreaker` patterns to prevent cascading failures if downstream Google services encounter latency.
* **Connection Pooling**: Optimized Redis buffers with a pool size of 100 and 20 minimum idle connections to handle massive telemetry surges.
* **Worker Pools**: Concurrency-safe telemetry processing using buffered channels and worker pools to prevent memory exhaustion (OOM events).

---

## ⚙️ How the Solution Works

1.  **Telemetry Ingestion**: IoT sensors or mobile applications send real-time coordinates and density data to the hardened `/telemetry` endpoint.
2.  **Validation**: The **Inbound Adapter** validates the data structure using `validator/v10` before it reaches the domain logic, ensuring "Always-Fail" defaults for missing properties.
3.  **Real-Time Processing**: The **Routing Service** analyzes current zone density against stadium capacity to determine congestion.
4.  **Buffer & Persistence**: High-frequency state is buffered in **Redis** for sub-millisecond heatmap generation, while long-term telemetry is persisted in **GCP Firestore**.
5.  **Intelligence Output**: The platform provides real-time "Reroute Events" to redirect fans to less congested gates or amenities.

---

## 🛡️ Security & Observability
* **Zero-Trust Containers**: Deployed using **Google Distroless (Non-Root)** images to eliminate shell access and minimize the attack surface.
* **Secret Management**: All credentials (API keys, DB credentials) are dynamically fetched from **GCP Secret Manager**; standard `.env` files are prohibited in production.
* **Probing**: Native `/healthz` (Liveness) and `/readyz` (Readiness) endpoints allow **Google Cloud Run** to monitor connectivity and self-heal the service.
* **Graceful Shutdown**: The system intercepts `SIGTERM` signals, ensuring all background worker routines complete their pipeline buffers before termination.
* **Vulnerability Management**: Automated **SecOps Sweep** via `govulncheck` in the CI/CD pipeline to block deployments with known vulnerabilities.

---

## 📝 Assumptions Made
* **Network Reliability**: It is assumed that the event venue provides sufficient local network infrastructure (5G/Wi-Fi) to transmit telemetry with sub-second latency.
* **Identity Management**: User identity is assumed to be handled by an external authentication provider, allowing this backend to focus on telemetry and routing logic.
* **GCP Permissions**: Deployment assumes the Cloud Run Service Account has been granted the `Secret Manager Secret Accessor` role.

---
© 2026 Aman Sayyad since 15 April 2006. All Rights Reserved.
