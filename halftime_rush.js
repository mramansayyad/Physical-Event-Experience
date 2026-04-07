import http from 'k6/http';
import { check, sleep } from 'k6';

// Halftime Rush Execution Plan: Target 50,000 concurrency in 120s
export const options = {
  stages: [
    { duration: '30s', target: 5000 },    // Gradual ramp-up
    { duration: '60s', target: 50000 },   // HALFTIME RUSH: The main spike
    { duration: '30s', target: 50000 },   // Sustain maximum crowd capacity
    { duration: '10s', target: 0 },       // Dissipate
  ],
  thresholds: {
    // Dynamic performance boundaries targeting sub-100ms
    http_req_duration: ['p(95)<100'], 
    http_req_failed: ['rate<0.01'],   // Error rate should be less than 1%
  },
};

export default function () {
  const payload = JSON.stringify({
    device_id: `fan-${__VU}-${__ITER}`,
    latitude: 34.0522 + (Math.random() * 0.001),
    longitude: -118.2437 + (Math.random() * 0.001),
    timestamp: new Date().toISOString()
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      // Simulation IAP JWT mapping
      'x-goog-iap-jwt-assertion': 'mock-valid-jwt-token'
    },
  };

  // Mock ingress triggering the telemetry pipeline boundary
  const res = http.post('https://stadium-experience.internal.run/telemetry/ingest', payload, params);

  check(res, {
    'is status 202': (r) => r.status === 202,
    'handled fast': (r) => r.timings.duration < 50
  });

  // Calculate random 1-2s delay between telemetry pings
  sleep(Math.random() * 2 + 1);
}
