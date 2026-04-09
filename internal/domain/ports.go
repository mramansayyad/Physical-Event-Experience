package domain

import "context"

// -- Input Ports (Invoked by Transports/Handlers) --

// RoutingUseCase handles dynamic crowd direction.
type RoutingUseCase interface {
	AnalyzeCongestion(ctx context.Context, userID, currentZone string, density float64) (*RerouteEvent, error)
}

// WaitTimeUseCase abstracts queue predictions.
type WaitTimeUseCase interface {
	PredictWaitTime(ctx context.Context, amenityID string) (int, error)
}

// -- Output Ports (Implemented by Adapters) --

// HeatmapReader retrieves read-only venue state.
type HeatmapReader interface {
	GetZoneHeatmap(ctx context.Context, zoneID string) (Heatmap, error)
	ListGateHeatmaps(ctx context.Context) ([]Heatmap, error)
}

// TelemetryWriter writes high-frequency stream events.
type TelemetryWriter interface {
	BufferTelemetry(ctx context.Context, record TelemetryRecord) error
}

// LocationRepository combines reads and writes for generalized state persistence.
type LocationRepository interface {
	UpdateUserLocation(ctx context.Context, userID string, loc Location) error
	GetZoneTelemetry(ctx context.Context, zoneID string) ([]TelemetryRecord, error)
}
