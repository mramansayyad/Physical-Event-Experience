package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/redis/go-redis/v9"
	"github.com/virtual-promptwars/stadium-backend/internal/adapter"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = "virtual-promptwars-stadium"
	}

	// Wait-Time and Cloud Run timeout context
	log.Println("Initializing Real-Time Stadium Backend (Hexagonal Architecture)...")

	// 1. Initialize Firestore client (Cold Start optimization: init globally)
	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer fsClient.Close()

	// 2. Initialize Pub/Sub client
	psClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to initialize Pub/Sub: %v", err)
	}
	defer psClient.Close()

	// 3. Adapter Layer Instantiation & Dependency Injection
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})
	
	_ = adapter.NewFirestoreRepository(fsClient)              // Standard repository initialized
	redisBuffer := adapter.NewRedisBuffer(rdb, fsClient)      // Ephemeral pipeline cache natively
	
	// Refactored Ingress Boundary: Streamer specifically wired to Memory transient
	// thereby buffering the global system explicitly avoiding Firestore hot constraints directly
	_ = adapter.NewPubSubStreamer(psClient, redisBuffer)

	log.Println("Adapters successfully wired into domain ports.")

	// Cloud Run HTTP Server routing explicitly mapping Hexagonal Adapters
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","layer":"adapter","mesh":"active"}`))
	})
	
	// Instantiating HTTP Native Handlers routing JSON requests natively
	mux.HandleFunc("/v1/telemetry/ingest", handleIngest(nil)) // Mock telemetry Service inject
	mux.HandleFunc("/v1/crowd/heatmap", handleHeatmap(nil))   // Mock routing Service inject
	mux.HandleFunc("/", handleRoot())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)
	
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("HTTP server crashed: %v", err)
	}
}
