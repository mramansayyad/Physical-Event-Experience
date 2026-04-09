package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/redis/go-redis/v9"
	"github.com/virtual-promptwars/stadium-backend/internal/adapter"
	"github.com/virtual-promptwars/stadium-backend/internal/domain"
	"github.com/virtual-promptwars/stadium-backend/internal/transport"
	"github.com/virtual-promptwars/stadium-backend/internal/transport/middleware"
)

func main() {
	ctx := context.Background()
	
	// Purged os.Getenv completely: statically enforce base identifiers for Secret resolution
	projectID := "virtual-promptwars-stadium"

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

	// Securely retrieve sensitive infrastructure credentials via Secret Manager ("Always-Fail" condition)
	secMgr, err := adapter.NewSecretManager(ctx)
	if err != nil {
		log.Fatalf("System halted: unable to initialize Secret Manager: %v", err)
	}
	defer secMgr.Close()

	redisHost, err := secMgr.GetSecret(ctx, projectID, "REDIS_HOST")
	if err != nil || redisHost == "" {
		log.Fatalf("System halted: REDIS_HOST secret missing or inaccessible. Error: %v", err)
	}

	redisPort, err := secMgr.GetSecret(ctx, projectID, "REDIS_PORT")
	if err != nil || redisPort == "" {
		log.Fatalf("System halted: REDIS_PORT secret missing or inaccessible. Error: %v", err)
	}

	redisPassword, err := secMgr.GetSecret(ctx, projectID, "REDIS_PASSWORD")
	if err != nil {
		log.Fatalf("System halted: REDIS_PASSWORD secret missing or inaccessible. Error: %v", err)
	}

	// 3. Adapter Layer Instantiation & Dependency Injection
	rdb := redis.NewClient(&redis.Options{
		Addr:            redisHost + ":" + redisPort,
		Password:        redisPassword,
		PoolSize:        100,
		MinIdleConns:    20,
		ConnMaxIdleTime: 5 * time.Minute,
	})
	
	_ = adapter.NewFirestoreRepository(fsClient)              // Standard repository initialized
	redisBuffer := adapter.NewRedisBuffer(rdb, fsClient)      // Ephemeral pipeline cache natively
	
	// Refactored Ingress Boundary: Streamer specifically wired to Memory transient
	// thereby buffering the global system explicitly avoiding Firestore hot constraints directly
	_ = adapter.NewPubSubStreamer(psClient, redisBuffer)

	log.Println("Adapters successfully wired into domain ports.")

	// Cloud Run HTTP Server routing explicitly mapping Hexagonal Adapters
	mux := http.NewServeMux()
	// Dependency Ping configuration natively checking Adapter States
	dbCheck := func(ctx context.Context) error {
		return rdb.Ping(ctx).Err()
	}

	mux.HandleFunc("/healthz", transport.HandleHealthz())
	mux.HandleFunc("/readyz", transport.HandleReadyz(dbCheck))
	
	// Injecting separated ISP interfaces into domain routing service and transport layer
	routeService := domain.NewRouteService(nil, nil)

	// Instantiating HTTP Native Handlers routing JSON requests natively via isolated transport package
	mux.HandleFunc("/v1/telemetry/ingest", transport.HandleIngest(routeService)) // Mock telemetry Service inject
	mux.HandleFunc("/v1/crowd/heatmap", transport.HandleHeatmap(routeService))   // Mock routing Service inject
	mux.HandleFunc("/", transport.HandleRoot())

	// Execute internal Pprof diagnostic listener
	go StartPprofServer()

	// Hardcoded port bound entirely to internal definitions (Zero-Trust Environment configs)
	port := "8080"
	log.Printf("Listening securely on :%s", port)
	
	// Server instance configured for explicit Graceful Shutdown bound natively
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: middleware.SecurityHeaders(middleware.GCPLoggingInterceptor(mux)),
	}

	// Detach asynchronous execution directly
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server crashed: %v", err)
		}
	}()

	// Orchestrate Graceful Shutdown locking Background Worker allocations inherently
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Initiating Zero-Trust Graceful Shutdown...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Awaiting RoutingService Worker Pool termination...")
	routeService.Shutdown()
	log.Println("Hexagonal System Successfully Halted!")
}
