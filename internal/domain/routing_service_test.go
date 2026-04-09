package domain

import (
	"context"
	"testing"
)

func TestAnalyzeCongestion_TableDriven(t *testing.T) {
	svc := NewRouteService(nil, nil)

	tests := []struct {
		name           string
		density        float64
		gates          []Heatmap
		expectedTarget string
		expectEvent    bool
	}{
		{
			name:    "Maximum Capacity Pings (>80%) - Targets Best Gate",
			density: 0.85,
			gates: []Heatmap{
				{ZoneID: "gate-a", DensityLevel: 0.65}, // Moderate
				{ZoneID: "gate-b", DensityLevel: 0.35}, // Low <40%, ideal target
				{ZoneID: "gate-c", DensityLevel: 0.95}, // Critical
			},
			expectedTarget: "gate-b",
			expectEvent:    true,
		},
		{
			name:    "Low Density (<80%) - No Event Needed",
			density: 0.50,
			gates:   []Heatmap{},
			expectedTarget: "",
			expectEvent:    false,
		},
		{
			name:    "Zero Density Edge Case - No Event",
			density: 0.00,
			gates: []Heatmap{
				{ZoneID: "gate-a", DensityLevel: 0.20},
			},
			expectedTarget: "",
			expectEvent:    false,
		},
		{
			name:    "Severe Congestion but NO Available Gates (<40%) - No Event",
			density: 0.90,
			gates: []Heatmap{
				{ZoneID: "gate-a", DensityLevel: 0.80},
				{ZoneID: "gate-b", DensityLevel: 0.50}, // Stalled
			},
			expectedTarget: "",
			expectEvent:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := svc.AnalyzeCongestion(context.Background(), "fan-123", "sector-1", tt.density, tt.gates)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			
			if tt.expectEvent {
				if event == nil {
					t.Fatalf("expected reroute event, got nil")
				}
				if event.TargetGateID != tt.expectedTarget {
					t.Errorf("expected target gate %s, got %s", tt.expectedTarget, event.TargetGateID)
				}
			} else {
				if event != nil {
					t.Fatalf("expected NO reroute event, got %v", event)
				}
			}
		})
	}
}

func TestPredictWaitTime_WeightedAverage(t *testing.T) {
	svc := NewRouteService(nil, nil)

	// Weighted avg for [10, 20]: (10*1 + 20*2) / 3 = 16.66 -> Int = 16
	samples := []int{10, 20}
	prediction := svc.PredictWaitTime(context.Background(), samples)
	
	if prediction != 16 {
		t.Errorf("Expected WaitTime of 16 mins, got %d", prediction)
	}
}
