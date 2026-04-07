#!/bin/bash
# Ensures Redis IPs never reside statically in variables locally via CI/CD pipelines
set -e

echo "Extracting dynamic Redis limits directly from TF State allocations natively..."

# Mapping terraform metrics silently avoiding native terminal log outputs
REDIS_HOST=$(terraform -chdir=../terraform output -raw redis_internal_ip)
REDIS_PORT=$(terraform -chdir=../terraform output -raw redis_port)

echo "Injecting extracted logic mapping seamlessly strictly into Google Cloud Secret Manager..."

echo -n $REDIS_HOST | gcloud secrets create stadium-redis-host --data-file=- --replication-policy=automatic
echo -n $REDIS_PORT | gcloud secrets create stadium-redis-port --data-file=- --replication-policy=automatic

# Lock specific runtime access explicitly allocating Cloud Run execution Service Account perfectly restricting unauthorized views
gcloud secrets add-iam-policy-binding stadium-redis-host \
    --member="serviceAccount:stadium-runner@stadium-experience-loc.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding stadium-redis-port \
    --member="serviceAccount:stadium-runner@stadium-experience-loc.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

echo "Variables securely injected masking plain-text states completely."
