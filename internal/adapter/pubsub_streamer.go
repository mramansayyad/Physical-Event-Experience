package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/virtual-promptwars/stadium-backend/internal/domain"
	
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type PubSubStreamer struct {
	client *pubsub.Client
	repo   domain.LocationRepository
}

func NewPubSubStreamer(client *pubsub.Client, repo domain.LocationRepository) *PubSubStreamer {
	return &PubSubStreamer{
		client: client,
		repo:   repo,
	}
}

// IngestTelemetry streams incoming telemetry from stadium sensors.
func (p *PubSubStreamer) IngestTelemetry(ctx context.Context, subscriptionID string) error {
	sub := p.client.Subscription(subscriptionID)
	
	err := sub.Receive(ctx, func(msgCtx context.Context, msg *pubsub.Message) {
		// INJECT OTEL TRACING
		msgCtx, span := otel.Tracer(tracerName).Start(msgCtx, "PubSub.Ingress")
		// The custom attributes required for the payload
		span.SetAttributes(attribute.String("message_id", msg.ID))
		defer span.End()

		var record domain.TelemetryRecord
		if err := json.Unmarshal(msg.Data, &record); err != nil {
			log.Printf("Failed to unmarshal telemetry message natively: %v", err)
			
			// Maps structurally non-compliant payload sequences accurately directly to Dead-Letter Queues inherently
			dlq := p.client.Topic("telemetry-dlq")
			// Asynchronous native publishing avoiding blocking main stream execution 
			dlq.Publish(msgCtx, &pubsub.Message{
				Data: msg.Data,
				Attributes: map[string]string{
					"error_reason": "malformed_json_violation",
					"original_id":  msg.ID,
				},
			})
			
			msg.Ack() // Complete receipt structurally bypassing endless retry pipeline blocks natively
			return
		}
		
		// Forward device telemetry to FireStore span mapping
		err := p.repo.UpdateUserLocation(msgCtx, record.DeviceID, record.Location)
		if err != nil {
			log.Printf("Failed to update user location: %v", err)
			msg.Nack()
			return
		}
		
		// Conclude logic map
		msg.Ack()
	})
	
	if err != nil && err != context.Canceled {
		return fmt.Errorf("pubsub stream terminated with error: %v", err)
	}
	
	return nil
}
