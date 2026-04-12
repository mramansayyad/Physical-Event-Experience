package domain

import (
	"errors"
	"time"
)

var (
	// ErrZoneNotFound is returned when a requested stadium zone does not exist.
	ErrZoneNotFound = errors.New("zone not found in stadium topology")
	// ErrTelemetryInvalid is returned when inbound telemetry is malformed or out of bounds.
	ErrTelemetryInvalid = errors.New("telemetry record violates domain constraints")
	// ErrRoutingFailed is returned when a routing request fails to find a valid pathway.
	ErrRoutingFailed = errors.New("failed to calculate optimal congestion routing")
)

// Heatmap represents the localized density of attendees dynamically operating in a specific zone.
type Heatmap struct {
	ZoneID       string
	DensityLevel float64 // 0.0 to 1.0 representing aggregated crowd density capacity
	Timestamp    time.Time
}

// Location represents rigorous geographical coordinates within the stadium boundaries.
type Location struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

// TelemetryRecord aggregates high-frequency data points for positional foot traffic ingestion organically.
type TelemetryRecord struct {
	DeviceID  string    `json:"device_id" validate:"required,uuid4"`
	Location  Location  `json:"location" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required,past_timestamp"`
}
