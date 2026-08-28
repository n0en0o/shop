import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const errorRate = new Rate('grpc_errors');
const promotionDuration = new Trend('promotion_grpc_duration', true);

const TARGET = __ENV.GRPC_TARGET || 'promotion-grpc:7003';
const client = new grpc.Client();
client.load(['/proto'], 'promo.proto');

export const options = {
  stages: [
    { duration: '10s', target: 10 },
    { duration: '20s', target: 100 },
    { duration: '20s', target: 50 },
    { duration: '1m', target: 300 },
    { duration: '5m', target: 100 },
  ],
  thresholds: {
    checks: ['rate>0.9'],
    grpc_errors: ['rate<0.1'],
    promotion_grpc_duration: ['p(95)<500', 'p(99)<1000'],
  },
};

const catalogItemIds = [
  '21212121-2121-2121-2121-212121212121',
  '20202020-2020-2020-2020-202020202020',
  '27272727-2727-2727-2727-272727272727',
];

function pickCatalogItemId() {
  return catalogItemIds[Math.floor(Math.random() * catalogItemIds.length)];
}

export default function () {
  client.connect(TARGET, { plaintext: true });

  const catalogItemId = pickCatalogItemId();
  const started = Date.now();
  const response = client.invoke('promotion.PromotionService/GetPromoByCatalogItem', {
    catalogItemId,
  });
  const duration = Date.now() - started;

  const ok = check(response, {
    'status OK': (r) => r && r.status === grpc.StatusOK,
    'promotion exists': (r) => r && r.message && r.message.promotion,
    'catalog item matches': (r) =>
      r && r.message && r.message.promotion && r.message.promotion.catalogItemId === catalogItemId,
    'response time < 500ms': () => duration < 500,
  });

  errorRate.add(!ok);
  promotionDuration.add(duration);

  client.close();
  sleep(0.1);
}
