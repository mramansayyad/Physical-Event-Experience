package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/redis/go-redis/v9"
	"github.com/virtual-promptwars/stadium-backend/internal/adapter"
	"github.com/virtual-promptwars/stadium-backend/internal/domain"
)

// Syncer Worker triggered locally via Cloud Scheduler -> Cloud Run Job natively every 5s.
func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Firestore init native error: %v", err)
	}
	defer fsClient.Close()
	
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})
	
	firestoreRepo := adapter.NewFirestoreRepository(fsClient)
	redisBuffer := adapter.NewRedisBuffer(rdb, fsClient)

	log.Println("Executing Syncer worker. Isolating ephemeral structures...")

	// Atomically pull all cached aggregates spanning instances and delete transient state directly.
	aggregates, err := redisBuffer.FlushAggregatedTotals(ctx)
	if err != nil {
		log.Fatalf("Locked pulling buffer logic: %v", err)
	}

	var heatmaps []domain.Heatmap
	for zone, countStr := range aggregates {
		count, _ := strconv.Atoi(countStr)
		// Transform ingress volume into relative percentage density against a static structural pool
		density := float64(count) / 1000.0 
		
		heatmaps = append(heatmaps, domain.Heatmap{
			ZoneID:       zone,
			DensityLevel: density,
			Timestamp:    time.Now(),
		})
	}

	if len(heatmaps) > 0 {
		// Single committed pipeline strictly updating persistent records via 1 bulk mapped atomic Firestore Batch
		err = firestoreRepo.BatchUpdateHeatmaps(ctx, heatmaps)
		if err != nil {
			log.Fatalf("Catastrophic Sync Native Failure: %v", err)
		}
		log.Printf("Successfully completed Memory-to-Storage boundary payload. %d unique zone aggregates mapped securely.", len(heatmaps))
	} else {
		log.Println("Zero ingress latency. Standing by natively.")
	}
}
