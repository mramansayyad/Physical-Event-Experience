# Real-Time Stadium Experience Platform

## Overview
A high-performance backend serving Dynamic Crowd Routing and Wait-Time Prediction features for stadium attendees. Built purely under strict Hexagonal Architecture constraints and deployed to Google Cloud Run instances via a Distroless footprint.

## JSON:API Documentation

### 1. Crowd Heatmap Endpoint
**`GET /v1/crowd/heatmap?zone={zone_id}`**

Retrieves the real-time density matrix of any specific structural section directly from the Firestore sync stream.

**Payload:**
```json
{
  "data": {
    "type": "heatmaps",
    "id": "H-12345",
    "attributes": {
      "zone_id": "section-104",
      "density_level": 0.85,
      "congestion_status": "HIGH",
      "timestamp": "2026-04-07T20:53:00Z"
    }
  }
}
```

### 2. Stalls Wait-Times Endpoint
**`GET /v1/stalls/wait-times?amenity={amenity_id}`**

Calculates projected wait durations via a localized Weighted Moving Average mapping over Pub/Sub telemetry ingress.

**Payload:**
```json
{
  "data": {
    "type": "wait-times",
    "id": "W-98765",
    "attributes": {
      "amenity_id": "concessions-gate-4",
      "wait_time_minutes": 14
    }
  }
}
```

## Securing Access (Identity-Aware Proxy / IAP)
The API strictly enforces a zero-trust security mesh layer using Google's IAP interface.
1. Deploy the API via Cloud Run directly without opening the public ingress trigger.
2. Setup a **Global HTTP(S) Load Balancer** with a Serverless Network Endpoint Group (NEG) terminating against the Cloud Run service.
3. Subnet the Global Load Balancer to a strict IP origin.
4. In GCC, navigate to **Security > Identity-Aware Proxy**. Turn on IAP for the newly routed Backend Service.
5. Assign `roles/iap.httpsResourceAccessor` to the target domain members.
6. The Backend adapters now uniquely extract the valid JSON Web Token via the `x-goog-iap-jwt-assertion` ingress headers on all requests.

## Scaling Architecture: MemoryStore (Redis) & VPC Connectors
To process 50,000+ concurrent telemetry ingestion spikes without breaching Firestore index write constraints, our architecture relies strictly on an Ephemeral Buffer:
1. Initialize a **Google Cloud Memorystore (Redis)** instance locally assigned bounding pure internal IP logic.
2. Provision a **Serverless VPC Access Connector** binding your Cloud Run functions directly to the Redis internal network isolation constraints.
3. Configure --vpc-egress=private-ranges-only inside Knative specifications minimizing outbound proxies securely bridging REDIS_HOST.
