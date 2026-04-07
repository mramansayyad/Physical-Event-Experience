package domain

import (
	"context"
	"testing"
)

// mockLocRepo is a stub adapter to validate Hexagonal dependency injection
type mockLocRepo struct{}

func (m *mockLocRepo) UpdateUserLocation(ctx context.Context, userID string, loc Location) error {
	return nil
}
func (m *mockLocRepo) GetZoneTelemetry(ctx context.Context, zoneID string) ([]TelemetryRecord, error) {
	return nil, nil
}

func TestAnalyzeCongestion_TriggersReroute(t *testing.T) {
	repo := &mockLocRepo{}
	svc := NewRouteService(repo)

	// High congestion density (>80%)
	userID := "fan-789"
	currentZone := "section-104"
	currentDensity := 0.85

	// Nearest possible gates
	gates := []Heatmap{
		{ZoneID: "gate-a", DensityLevel: 0.65}, // Moderate
		{ZoneID: "gate-b", DensityLevel: 0.35}, // Low <40%, ideal target
		{ZoneID: "gate-c", DensityLevel: 0.95}, // Critical
	}

	event, err := svc.AnalyzeCongestion(context.Background(), userID, currentZone, currentDensity, gates)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if event == nil {
		t.Fatal("Expected a RerouteEvent to be generated, got nil")
	}
	if event.TargetGateID != "gate-b" {
		t.Errorf("Expected optimal route to gate-b, got %s", event.TargetGateID)
	}
}

func TestAnalyzeCongestion_NoReroute(t *testing.T) {
	repo := &mockLocRepo{}
	svc := NewRouteService(repo)

	// Density below the 80% threshold
	event, err := svc.AnalyzeCongestion(context.Background(), "fan-123", "section-105", 0.50, []Heatmap{})
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if event != nil {
		t.Fatal("Expected no reroute event, got an event")
	}
}

func TestPredictWaitTime_WeightedAverage(t *testing.T) {
	repo := &mockLocRepo{}
	svc := NewRouteService(repo)

	// Weighted avg for [10, 20]: (10*1 + 20*2) / 3 = 16.66 -> Int = 16
	samples := []int{10, 20}
	prediction := svc.PredictWaitTime(context.Background(), samples)
	
	if prediction != 16 {
		t.Errorf("Expected WaitTime of 16 mins, got %d", prediction)
	}
}
