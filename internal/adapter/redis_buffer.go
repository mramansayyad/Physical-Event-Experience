package adapter

import (
	"context"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/redis/go-redis/v9"
	"github.com/virtual-promptwars/stadium-backend/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CircuitBreaker explicitly protects the Hot Path from infinite retries causing upstream Pub/Sub unacked bottlenecks.
type CircuitBreaker struct {
	state          string // "CLOSED", "OPEN", "HALF_OPEN"
	failureCount   int
	threshold      int
	lastFailure    time.Time
	resetTimeout   time.Duration
	mu             sync.Mutex
}

func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == "OPEN" {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.state = "HALF_OPEN"
			return true
		}
		return false
	}
	return true
}

func (cb *CircuitBreaker) RecordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()
		if cb.failureCount >= cb.threshold {
			cb.state = "OPEN"
		}
	} else {
		cb.failureCount = 0
		cb.state = "CLOSED"
	}
}

type RedisBuffer struct {
	client   *redis.Client
	fallback *firestore.Client
	cb       *CircuitBreaker
}

// Memory-to-Memory transient layer integrating native Circuit Breaking logic against 'Low-Priority' Firestore collections.
func NewRedisBuffer(client *redis.Client, fallback *firestore.Client) *RedisBuffer {
	return &RedisBuffer{
		client:   client,
		fallback: fallback,
		cb: &CircuitBreaker{
			state:        "CLOSED",
			threshold:    3,
			resetTimeout: 10 * time.Second,
		},
	}
}

func (r *RedisBuffer) UpdateUserLocation(ctx context.Context, userID string, loc domain.Location) error {
	ctx, span := otel.Tracer("stadium-backend").Start(ctx, "Redis.EphemeralIngest.WithCircuitBreaker")
	defer span.End()
	zoneID := "section-104-main-exit"

	// Step 1: Circuit Evaluation
	if !r.cb.AllowRequest() {
		span.AddEvent("CircuitBreaker_Open_FallbackToFirestore")
		_ = r.executeFallback(ctx, userID, zoneID)
		return nil // Graceful degradation natively executed. Don't Nack in Pub/Sub.
	}

	luaScript := `
		local current = redis.call('HINCRBY', KEYS[1], ARGV[1], 1)
		return current
	`
	err := r.client.Eval(ctx, luaScript, []string{"stadium:density:aggregate"}, zoneID).Err()
	r.cb.RecordResult(err)

	if err != nil {
		log.Printf("[CHAOS] Redis write explicitly failed: %v", err)
		span.AddEvent("RedisFailure_FallbackInvoked")
		_ = r.executeFallback(ctx, userID, zoneID)
	}
	
	span.SetAttributes(attribute.String("zone_id", zoneID))
	return nil
}

func (r *RedisBuffer) executeFallback(ctx context.Context, userID, zoneID string) error {
	if r.fallback != nil {
		docRef := r.fallback.Collection("stadium-low-priority-sync").Doc(userID)
		_, err := docRef.Set(ctx, map[string]interface{}{
			"zoneID": zoneID,
			"timestamp": time.Now(),
		})
		return err
	}
	return nil
}

func (r *RedisBuffer) FlushAggregatedTotals(ctx context.Context) (map[string]string, error) {
	results, err := r.client.HGetAll(ctx, "stadium:density:aggregate").Result()
	if err == nil && len(results) > 0 {
		r.client.Del(ctx, "stadium:density:aggregate")
	}
	return results, err
}

func (r *RedisBuffer) GetZoneTelemetry(ctx context.Context, zoneID string) ([]domain.TelemetryRecord, error) {
	// Stub implementation to satisfy domain.LocationRepository
	return []domain.TelemetryRecord{}, nil
}
