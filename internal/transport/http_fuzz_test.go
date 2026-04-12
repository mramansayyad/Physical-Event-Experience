package transport

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// FuzzIngest applies structural mutation testing verifying native parsing constraints
// enforcing Memory safety bounds natively by bombarding the HTTP Ingest with bad JSON.
func FuzzIngest(f *testing.F) {
	// Native structurally valid seed preventing naive early-rejection blocks
	f.Add([]byte(`{"device_id":"550e8400-e29b-41d4-a716-446655440000","location":{"latitude":45.123,"longitude":-122.456},"timestamp":"2023-01-01T12:00:00Z"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"device_id": "missing_location"}`))
	f.Add([]byte(`[array_not_object]`))
	f.Add([]byte(`{"device_id":"550e8400-e29b-41d4-a716-446655440000","location":{"latitude":9000,"longitude":-180},"timestamp":"2024-01-01T00:00:00Z"}`)) // Bad Latitude mapping

	// Isolate execution bounds
	handler := HandleIngest(nil)

	f.Fuzz(func(t *testing.T, payload []byte) {
		req, err := http.NewRequest(http.MethodPost, "/v1/telemetry/ingest", bytes.NewReader(payload))
		if err != nil {
			return
		}
		
		recorder := httptest.NewRecorder()
		
		// If HandleIngest panics upon mutated byte sequences natively,
		// standard Go Fuzzing logic will organically mark execution as Failed, trapping the error state natively.
		handler(recorder, req)
		
		// Ensure standard protocol error boundaries natively map securely avoiding 500s inherently
		if recorder.Code == http.StatusInternalServerError {
			t.Errorf("Ingress boundary failed safely blocking JSON mutation: 500 status returned inherently.")
		}
	})
}
