package adapter

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/virtual-promptwars/stadium-backend/internal/domain"
	"google.golang.org/api/iterator"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "stadium-backend/adapter"

type FirestoreRepository struct {
	client *firestore.Client
}

func NewFirestoreRepository(client *firestore.Client) *FirestoreRepository {
	return &FirestoreRepository{client: client}
}

// UpdateUserLocation saves the latest coordinates.
func (r *FirestoreRepository) UpdateUserLocation(ctx context.Context, userID string, loc domain.Location) error {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "Firestore.UpdateUserLocation")
	span.SetAttributes(attribute.String("user_id", userID))
	defer span.End()

	docRef := r.client.Collection("locations").Doc(userID)
	
	_, err := docRef.Set(ctx, map[string]interface{}{
		"Latitude":  loc.Latitude,
		"Longitude": loc.Longitude,
		"Timestamp": time.Now(),
		"TTL":       time.Now().Add(5 * time.Minute), 
	}, firestore.MergeAll)
	
	return err
}

// GetZoneTelemetry retrieves telemetry records.
func (r *FirestoreRepository) GetZoneTelemetry(ctx context.Context, zoneID string) ([]domain.TelemetryRecord, error) {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "Firestore.GetZoneTelemetry")
	span.SetAttributes(attribute.String("zone_id", zoneID))
	defer span.End()

	var records []domain.TelemetryRecord
	staleThreshold := time.Now().Add(-5 * time.Minute)
	
	iter := r.client.Collection("telemetry").
		Where("ZoneID", "==", zoneID).
		Where("Timestamp", ">=", staleThreshold).
		Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate telemetry docs: %v", err)
		}
		
		var raw struct {
			DeviceID  string    `firestore:"DeviceID"`
			Latitude  float64   `firestore:"Latitude"`
			Longitude float64   `firestore:"Longitude"`
			Timestamp time.Time `firestore:"Timestamp"`
		}
		if err := doc.DataTo(&raw); err != nil {
			continue
		}
		
		records = append(records, domain.TelemetryRecord{
			DeviceID: raw.DeviceID,
			Location: domain.Location{
				Latitude:  raw.Latitude,
				Longitude: raw.Longitude,
			},
			Timestamp: raw.Timestamp,
		})
	}
	
	return records, nil
}

// BatchUpdateHeatmaps uses nested spans and batch writes.
func (r *FirestoreRepository) BatchUpdateHeatmaps(ctx context.Context, heatmaps []domain.Heatmap) error {
	ctx, span := otel.Tracer(tracerName).Start(ctx, "Firestore.BatchUpdateHeatmaps")
	span.SetAttributes(attribute.Int("batch_size", len(heatmaps)))
	defer span.End()

	batch := r.client.Batch()
	
	for _, hm := range heatmaps {
		docRef := r.client.Collection("heatmaps").Doc(hm.ZoneID)
		batch.Set(docRef, map[string]interface{}{
			"DensityLevel": hm.DensityLevel,
			"Timestamp":    hm.Timestamp,
			"TTL":          time.Now().Add(10 * time.Minute),
		}, firestore.MergeAll)
		
		// Record exact zone updates into OTEL for the specific Heatmap
		span.AddEvent("HeatmapProcessed", trace.WithAttributes(
			attribute.String("zone_id", hm.ZoneID),
			attribute.Float64("density", hm.DensityLevel),
		))
	}
	
	_, err := batch.Commit(ctx)
	return err
}
