import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const errorRate = new Rate('errors');
const catalogDuration = new Trend('catalog_duration', true);

const BASE_URL = __ENV.TARGET_URL || 'http://catalog-api:7001';

export const options = {
  stages: [
    { duration: '10s', target: 10 },
    { duration: '20s', target: 100 },
    { duration: '20s', target: 50 },
    { duration: '1m', target: 300 },
    { duration: '5m', target: 100 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    errors: ['rate<0.1'],
    checks: ['rate>0.9'],
    catalog_duration: ['p(95)<500'],
  },
};

const endpoints = [
  { url: '/api/v1/catalog-items', weight: 40 },
  { url: '/api/v1/brands', weight: 20 },
  { url: '/api/v1/categories', weight: 20 },
  { url: '/api/v2/catalog-items', weight: 15 },
  { url: '/health', weight: 5 },
];

function pickEndpoint() {
  const rand = Math.random() * 100;
  let cumulative = 0;

  for (const ep of endpoints) {
    cumulative += ep.weight;
    if (rand < cumulative) return ep.url;
  }

  return endpoints[0].url;
}

export default function () {
  const path = pickEndpoint();
  const res = http.get(`${BASE_URL}${path}`, {
    tags: { endpoint: path },
  });

  const ok = check(res, {
    'status 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  errorRate.add(!ok);
  catalogDuration.add(res.timings.duration);

  sleep(0.1);
}
