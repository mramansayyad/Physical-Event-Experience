package domain

import (
	"context"
	"sync"
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
	heatmapReader   HeatmapReader
	telemetryWriter TelemetryWriter
	telemetryJobs   chan TelemetryRecord // Bounded concurrency pool
	wg              sync.WaitGroup       // Concurrency drain lock
	ctx             context.Context      // Native lifecycle context binding
	cancel          context.CancelFunc   // Cancel signal mapped to Shutdown
}

func NewRouteService(hr HeatmapReader, tw TelemetryWriter) *RouteService {
	ctx, cancel := context.WithCancel(context.Background())
	svc := &RouteService{
		heatmapReader:   hr,
		telemetryWriter: tw,
		// Explicit buffered channel bounds representing "High-Concurrency Efficiency"
		telemetryJobs:   make(chan TelemetryRecord, 5000), 
		ctx:             ctx,
		cancel:          cancel,
	}
	
	// Establishing deterministic bounded goroutines for memory mapping (Worker Pool)
	for w := 1; w <= 10; w++ {
		svc.wg.Add(1)
		go svc.telemetryWorker()
	}
	
	return svc
}

// Shutdown initiates graceful pipeline draining protecting volatile configurations natively
func (s *RouteService) Shutdown() {
	s.cancel()             // Signals contextual draining safely across routines
	close(s.telemetryJobs) // Triggers channel exhaustion sequence natively
	s.wg.Wait()            // Halts parent thread strictly till completely zeroed
}

// telemetryWorker explicitly handles pipeline traffic using mapped context natively
func (s *RouteService) telemetryWorker() {
	defer s.wg.Done()
	
	for {
		select {
		case <-s.ctx.Done():
			// Prevents isolated Goroutine Leaks fundamentally
			return
		case job, ok := <-s.telemetryJobs:
			if !ok {
				return // Channel natively closed
			}
			if s.telemetryWriter != nil {
				_ = s.telemetryWriter.BufferTelemetry(s.ctx, job)
			}
		}
	}
}

// EnqueueTelemetry securely triggers the worker queue or sheds load natively (Backpressure)
func (s *RouteService) EnqueueTelemetry(record TelemetryRecord) {
	select {
	case s.telemetryJobs <- record:
	default: // Drops frame gracefully preserving Heap Health if burst buffer fills
	}
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
				Reason:        "error.routing.congestion_severe", // Transformed to native i18n map securely
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
