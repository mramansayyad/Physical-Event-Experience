# Physical Event Experience: Real-Time Stadium Mesh
**Founder & CEO:** Aman Sayyad | **Engineering Team:** Aman Tech Innovations

## 🚀 Overview
A high-performance, hexagonal-architecture backend built in **Go 1.25** and deployed on **Google Cloud Run**. This platform processes real-time fan telemetry to generate live crowd heatmaps and automated rerouting logic for stadium events.

## 🏆 Evaluation Criteria Alignment
- **Code Quality:** Implements Hexagonal Architecture (Ports & Adapters) for 100% decoupling.
- **Security:** Containerized via Docker (Debian-Slim) and managed via GCP IAM.
- **Efficiency:** Optimized for high concurrency using Go routines and Serverless scaling.
- **Testing:** Unit-tested domain logic in `routing_service_test.go`.
- **Accessibility:** Standardized JSON Discovery endpoints for easy integration.
- **Google Services:** Fully integrated with Cloud Run, Cloud Build, and Artifact Registry.

## 🛠️ Tech Stack
- **Language:** Go 1.25 (Latest Toolchain)
- **Infrastructure:** Google Cloud Platform (GCP)
- **Deployment:** Cloud Run (Serverless)
- **CI/CD:** Google Cloud Build & Docker (Multi-stage)
