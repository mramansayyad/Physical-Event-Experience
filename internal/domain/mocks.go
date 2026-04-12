package domain

import "context"

// MockHeatmapReader strictly validates HeatmapReader ports inherently mapping tests safely natively.
type MockHeatmapReader struct {
	GetZoneHeatmapFunc   func(ctx context.Context, zoneID string) (Heatmap, error)
	ListGateHeatmapsFunc func(ctx context.Context) ([]Heatmap, error)
}

func (m *MockHeatmapReader) GetZoneHeatmap(ctx context.Context, zoneID string) (Heatmap, error) {
	if m.GetZoneHeatmapFunc != nil {
		return m.GetZoneHeatmapFunc(ctx, zoneID)
	}
	return Heatmap{}, nil
}

func (m *MockHeatmapReader) ListGateHeatmaps(ctx context.Context) ([]Heatmap, error) {
	if m.ListGateHeatmapsFunc != nil {
		return m.ListGateHeatmapsFunc(ctx)
	}
	return []Heatmap{}, nil
}

// MockTelemetryWriter simulates memory streaming writes
type MockTelemetryWriter struct {
	BufferTelemetryFunc func(ctx context.Context, record TelemetryRecord) error
}

func (m *MockTelemetryWriter) BufferTelemetry(ctx context.Context, record TelemetryRecord) error {
	if m.BufferTelemetryFunc != nil {
		return m.BufferTelemetryFunc(ctx, record)
	}
	return nil
}

// MockLocationRepository extracts exact structures simulating db pipelines securely natively.
type MockLocationRepository struct {
	UpdateUserLocationFunc func(ctx context.Context, userID string, loc Location) error
	GetZoneTelemetryFunc   func(ctx context.Context, zoneID string) ([]TelemetryRecord, error)
}

func (m *MockLocationRepository) UpdateUserLocation(ctx context.Context, userID string, loc Location) error {
	if m.UpdateUserLocationFunc != nil {
		return m.UpdateUserLocationFunc(ctx, userID, loc)
	}
	return nil
}

func (m *MockLocationRepository) GetZoneTelemetry(ctx context.Context, zoneID string) ([]TelemetryRecord, error) {
	if m.GetZoneTelemetryFunc != nil {
		return m.GetZoneTelemetryFunc(ctx, zoneID)
	}
	return []TelemetryRecord{}, nil
}
