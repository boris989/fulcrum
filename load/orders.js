import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    scenarios: {
        ramping: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '30s', target: 10 },
                { duration: '1m',  target: 50 },
                { duration: '1m',  target: 100 },
                { duration: '30s', target: 0 },
            ],
            gracefulRampDown: '10s',
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.01'],     // <1% ошибок
        http_req_duration: ['p(95)<200'],   // p95 < 200ms (пример SLO)
    },
};

export default function () {
    const payload = JSON.stringify({
        user_id: 'load-user',
        amount: 100,
        currency: 'RUB',
    });

    const params = {
        headers: { 'Content-Type': 'application/json' },
    };

    const res = http.post('http://localhost:8080/orders', payload, params);

    check(res, {
        'status is 2xx': (r) => r.status >= 200 && r.status < 300,
    });

    sleep(0.1);
}