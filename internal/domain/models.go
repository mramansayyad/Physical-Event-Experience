package domain

import "time"

// Heatmap represents the density of attendees in a specific zone.
type Heatmap struct {
	ZoneID       string
	DensityLevel float64 // 0.0 to 1.0 representing crowd density
	Timestamp    time.Time
}

// Location represents geographical coordinates within the stadium.
type Location struct {
	Latitude  float64
	Longitude float64
}

// TelemetryRecord holds individual data points for foot traffic ingestion.
type TelemetryRecord struct {
	DeviceID  string
	Location  Location
	Timestamp time.Time
}
