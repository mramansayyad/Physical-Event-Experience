# Contributing to the Stadium Architecture

Welcome to the backend architecture for the Real-Time Stadium Experience Platform! This system guarantees sub-100ms processing for 50,000+ concurrent fan interactions bounding zero architectural lock-in natively.

## The Rule of Dependency Matrix (Hexagonal Architecture)

We strictly use **Ports and Adapters** (Hexagonal Architecture) to isolate and systematically protect the core business domains natively against fluctuating cloud mechanics and SDK updates dynamically.

### 1. Domain Layer (`internal/domain`)
- **Pure Go Logic**: This folder contains structural arrays and isolated logical boundaries modeling stadium entities (`Heatmap`, `WaitPrediction`). 
- **ABSOLUTE RESTRICTION**: You may **NOT** import any Google Cloud SDK (`cloud.google.com/go/*`), web server framework (`net/http`), or external database driver in this folder. Doing so critically breaks the framework natively and compromises local execution testing environments.
- **Ports**: Interactions natively mapping the external world are strictly structurally declared as simple Go interfaces here (e.g., `LocationRepository`).

### 2. Adapter Layer (`internal/adapter`)
- **Dirty Drivers**: Need to handle Firestore caching? Need to parse JSON from Pub/Sub streams? Connect to an external gRPC system? Do it explicitly natively here.
- **Injecting Interfaces**: Adapter structs must cleanly execute implementing the Ports systematically declared dynamically in the Domain natively exposing no infra structures. 
- **Redis vs Firestore**: The current `PubSubStreamer` intercepts incoming logic natively dumping directly down into Ephemeral Memory arrays (`RedisBuffer`). The independent `syncer` daemon handles migrating natively bounding persistent states smoothly back exactly into Firestore logic arrays natively.

## Verification
Breaking this architectural structure will efficiently fail the `audit_code.md` protocol actively blocking the CI/CD pipeline arrays natively ensuring zero corruption down the long-term execution natively. Keep the core logic precisely pure natively.
