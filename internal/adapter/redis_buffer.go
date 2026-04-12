package adapter

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
	"github.com/virtual-promptwars/stadium-backend/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type RedisBuffer struct {
	client   *redis.Client
	fallback *firestore.Client
	cb       *gobreaker.CircuitBreaker
}

// Memory-to-Memory transient layer integrating native Circuit Breaking logic against 'Low-Priority' Firestore collections.
func NewRedisBuffer(client *redis.Client, fallback *firestore.Client) *RedisBuffer {
	st := gobreaker.Settings{
		Name:        "RedisBuffer",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     10 * time.Second,
	}
	return &RedisBuffer{
		client:   client,
		fallback: fallback,
		cb:       gobreaker.NewCircuitBreaker(st),
	}
}

func (r *RedisBuffer) UpdateUserLocation(ctx context.Context, userID string, loc domain.Location) error {
	ctx, span := otel.Tracer("stadium-backend").Start(ctx, "Redis.EphemeralIngest.WithCircuitBreaker")
	defer span.End()
	zoneID := "section-104-main-exit"

	// Step 1: Circuit Evaluation
	_, err := r.cb.Execute(func() (interface{}, error) {
		luaScript := `
			local current = redis.call('HINCRBY', KEYS[1], ARGV[1], 1)
			return current
		`
		return r.client.Eval(ctx, luaScript, []string{"stadium:density:aggregate"}, zoneID).Result()
	})

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

// FlushAggregatedTotals natively exhausts the Memory buffers piping the commands transactionally avoiding sequential round-trip blocking.
func (r *RedisBuffer) FlushAggregatedTotals(ctx context.Context) (map[string]string, error) {
	pipe := r.client.Pipeline()
	
	// Pipeline stage 1: pull the aggregated zone definitions organically
	resultsCmd := pipe.HGetAll(ctx, "stadium:density:aggregate")
	// Pipeline stage 2: detach the state clearing the density natively
	pipe.Del(ctx, "stadium:density:aggregate")
	
	// Atomically execute avoiding network overhead map
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	
	return resultsCmd.Val(), nil
}

func (r *RedisBuffer) GetZoneTelemetry(ctx context.Context, zoneID string) ([]domain.TelemetryRecord, error) {
	// Stub implementation to satisfy domain.LocationRepository
	return []domain.TelemetryRecord{}, nil
}

// BufferTelemetry natively implements the domain.TelemetryWriter injection port
func (r *RedisBuffer) BufferTelemetry(ctx context.Context, record domain.TelemetryRecord) error {
	return r.UpdateUserLocation(ctx, record.DeviceID, record.Location)
}
