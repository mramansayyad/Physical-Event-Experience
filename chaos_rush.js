import http from 'k6/http';
import { check, sleep } from 'k6';

// Chaos Rush Execution Plan: Target 50,000 concurrency in 120s
export const options = {
  stages: [
    { duration: '30s', target: 5000 },    // Gradual ramp-up
    { duration: '60s', target: 50000 },   // HALFTIME RUSH: The main spike
    { duration: '30s', target: 50000 },   // Sustain maximum crowd capacity
    { duration: '10s', target: 0 },       // Dissipate
  ],
  thresholds: {
    // Expect 20% failures exactly because of the 10% JSON / 10% out of bounds injections
    http_req_failed: ['rate>=0.15'],
  },
};

export default function () {
  const isBrokenJson = Math.random() < 0.10;
  const isOutOfBounds = Math.random() < 0.10 && !isBrokenJson;

  let payload;
  let params = {
    headers: {
      'Content-Type': 'application/json',
      'x-goog-iap-jwt-assertion': 'mock-valid-jwt-token'
    },
  };

  if (isBrokenJson) {
    // 10% malformed JSON explicitly missing a struct quote mapped for Native unmarshal triggers
    payload = `{"device_id": "fan-${__VU}-${__ITER}", "latitude": 34.0522, longitude: -118.2437}`;
  } else if (isOutOfBounds) {
    // 10% Coordinate bounding exception mapped (Latitude > 90)
    payload = JSON.stringify({
      device_id: `fan-${__VU}-${__ITER}`,
      latitude: 150.00, // Invalid bounds mapped out
      longitude: -118.2437,
      timestamp: new Date().toISOString()
    });
  } else {
    payload = JSON.stringify({
      device_id: `fan-${__VU}-${__ITER}`,
      latitude: 34.0522 + (Math.random() * 0.001),
      longitude: -118.2437 + (Math.random() * 0.001),
      timestamp: new Date().toISOString()
    });
  }

  // Mock ingress triggering the telemetry pipeline boundary natively
  const res = http.post('http://127.0.0.1:8080/v1/telemetry/ingest', payload, params);

  // We should see a 400 bad request mapped inherently when garbage data is generated.
  if (isBrokenJson || isOutOfBounds) {
     check(res, {
        'is status 400 Bad Request': (r) => r.status === 400,
        'has RFC 7807 type': (r) => r.json("type") === "about:blank"
     });
  } else {
     check(res, {
        'is status 202 Accepted': (r) => r.status === 202,
     });
  }

  sleep(Math.random() * 2 + 1);
}
