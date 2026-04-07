package domain

import (
	"context"
)

// RerouteEvent encapsulates instructions for an attendee to avoid congestion.
type RerouteEvent struct {
	UserID        string
	CurrentZoneID string
	TargetGateID  string
	Reason        string
}

// RouteService implements routing and prediction business logic.
// Strictly adheres to Hexagonal Architecture (no infra imports).
type RouteService struct {
	locRepo LocationRepository
}

func NewRouteService(repo LocationRepository) *RouteService {
	return &RouteService{locRepo: repo}
}

// AnalyzeCongestion checks if a zone exceeds 80% capacity. 
// If so, it calculates the nearest available gate with <40% density and generates a RerouteEvent.
func (s *RouteService) AnalyzeCongestion(ctx context.Context, userID string, currentZone string, density float64, gates []Heatmap) (*RerouteEvent, error) {
	if density > 0.80 {
		var bestGate Heatmap
		found := false
		
		for _, gate := range gates {
			// Find gates below 40% density capacity
			if gate.DensityLevel < 0.40 {
				if !found || gate.DensityLevel < bestGate.DensityLevel {
					bestGate = gate
					found = true
				}
			}
		}
		
		if found {
			return &RerouteEvent{
				UserID:        userID,
				CurrentZoneID: currentZone,
				TargetGateID:  bestGate.ZoneID,
				Reason:        "Severe Congestion Avoidance",
			}, nil
		}
	}
	return nil, nil // No rerouting necessary
}

// PredictWaitTime uses a weighted moving average algorithm to predict wait times using recent samples.
// More recent samples (end of the slice) are weighted higher in the calculation.
func (s *RouteService) PredictWaitTime(ctx context.Context, recentSamples []int) int {
	if len(recentSamples) == 0 {
		return 0
	}
	
	var sumPriorities, sumWeights float64
	for i, sample := range recentSamples {
		weight := float64(i + 1) // Prioritize later (more recent) indices
		sumPriorities += float64(sample) * weight
		sumWeights += weight
	}
	
	// Truncate to nearest minute duration
	return int(sumPriorities / sumWeights)
}
