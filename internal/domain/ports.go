// Package domain contains the core business logic and models.
// It must have ZERO dependencies on external infrastructure (e.g., Firestore, Gin).

package domain

import (
	"context"
)

// CrowdScanner calculates and provides real-time heatmaps for dynamic crowd routing.
type CrowdScanner interface {
	// GenerateHeatmap calculates a heatmap delta for the given zone.
	GenerateHeatmap(ctx context.Context, zoneID string) (Heatmap, error)
	// AlertCongestion checks congestion thresholds and triggers reroute commands if necessary.
	AlertCongestion(ctx context.Context, heatmap Heatmap) error
}

// WaitTimePredictor predicts wait durations for stadium amenities (restrooms, concession stands).
type WaitTimePredictor interface {
	// PredictWaitTime returns the estimated wait time in minutes for a specific amenity.
	PredictWaitTime(ctx context.Context, amenityID string) (int, error)
}

// LocationRepository handles the retrieval and storage of geospatial and location data.
type LocationRepository interface {
	// UpdateUserLocation saves the latest coordinates for an attendee.
	UpdateUserLocation(ctx context.Context, userID string, loc Location) error
	// GetZoneTelemetry retrieves current telemetry records for a given zone.
	GetZoneTelemetry(ctx context.Context, zoneID string) ([]TelemetryRecord, error)
}
