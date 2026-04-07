# Disaster Recovery (DR) Executable: "The Big Red Button"

## Incident Scope: Permanent Redis/VPC Vaporization
If the MemoryStore VPC connection collapses completely (or the network isolated peering degrades), the internal Circuit Breaker will lock "OPEN". Concurrently, all dynamic ingest logic strictly bypasses the standard aggregation pipeline directly depositing raw positional arrays into the `stadium-low-priority-sync` unoptimized collection.

## Phase 1: Activate Direct-to-Firestore Override
When `alerts.yaml` triggers the PagerDuty P1 signaling permanent buffer stagnation across the backend, execute the explicit Network Override bridging the fallback directly across Version 1 architecture dynamically:

1. **Environmental Swap Protocol**:
   Natively set `DIRECT_FIRESTORE_OVERRIDE=true` in the main routing module.
   ```bash
   gcloud run services update stadium-experience-backend \
     --update-env-vars DIRECT_FIRESTORE_OVERRIDE=true \
     --region us-central1 \
     --project stadium-experience-loc
   ```

2. **Algorithm Reaction**:
   Upon detecting this flag natively, the `cmd/api/main.go` container explicitly zeroes out `RedisBuffer` dynamically. It forcefully injects `FirestoreRepository` isolating the logic perfectly back mapping Version 1 structurally.

3. **Rate Limit Client Apps**:
   To prevent 50,000 active fans globally crashing the Firestore logic maps without the transient Redis buffer, change standard polling arrays. Natively downgrade API limits dynamically inside `client_sdk.ts` modifying connection intervals from `5s` -> `60s`.

## Phase 2: Restoring Nominal Systems Array
1. Establish clean VPC connectivity natively parsing healthy redis pings against `/redis-healthz`.
2. Delete the override vector explicitly: `--remove-env-vars DIRECT_FIRESTORE_OVERRIDE`.
3. Wipe localized transient cache instances mapping `stadium-low-priority-sync`.
