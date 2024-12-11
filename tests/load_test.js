import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 1000, // Number of virtual users
    duration: '30m', // Duration of the test
};

export default function () {
    const url = 'http://localhost:8080/api/scrape';
    const payload = JSON.stringify({
        urls: [
            "https://map.naver.com/",
            "https://shopping.naver.com/",
            "https://shoppinglive.naver.com/",
        ],
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        'is status 200': (r) => r.status === 200,
        'is valid JSON': (r) => r.json('results') !== null,
    });

    sleep(1); // Simulate interval between requests
}
