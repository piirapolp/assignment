// Load testing script using K6 for assignment services
import http from 'k6/http';
import {check, sleep} from 'k6';
import {Counter, Rate, Trend} from 'k6/metrics';
import {SharedArray} from 'k6/data';

import {htmlReport} from 'https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js';
import {textSummary} from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';

// Custom metrics to measure performance

// api metrics
const getUserResponseTime = new Trend('get_user_response_time');
const loginResponseTime = new Trend('login_response_time');
const getUserAccountsResponseTime = new Trend('get_user_accounts_response_time');
const getUseDebitCardsResponseTime = new Trend('get_user_debit_cards_response_time');
const getUserSavedAccountsResponseTime = new Trend('get_user_saved_accounts_response_time');
const getUserBannersResponseTime = new Trend('get_user_banners_response_time');

// counter metrics
const errorRate = new Rate('error_rate');
const successfulGetUser = new Counter('successful_get_user');
const successfulLogin = new Counter('successful_login');
const successfulGetUserAccounts = new Counter('successful_get_user_accounts');
const successfulGetUserDebitCards = new Counter('successful_get_user_debit_cards');
const successfulGetUserSavedAccounts = new Counter('successful_get_user_saved_accounts');
const successfulGetUserBanners = new Counter('successful_get_user_banners');

const transactionCounter = new Counter('transaction_counter');

const BASE_URL = 'http://localhost:3000';

// Configuration: Max 500 concurrent connections on Fiber and DB
// To stay within safe limits, we'll target up to 400 VUs in test
export const options = {
    scenarios: {
        light_load: {
            executor: 'constant-vus',
            vus: 50,
            duration: '1m',
            startTime: "0s",
            tags: {scenario: 'light_load'},
        },
        normal_load: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '1m', target: 100 },
                { duration: '2m', target: 200 },
                { duration: '1m', target: 0 },
            ],
            startTime: "1m5s",
            tags: { scenario: 'normal_load' },
        },
        heavy_load: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '2m', target: 200 },
                { duration: '3m', target: 400 },
                { duration: '2m', target: 600 },
                { duration: '1m', target: 0 },
            ],
            startTime: "5m15s",
            tags: { scenario: 'heavy_load' },
        },
    },

    thresholds: {
        // Light Load
        'login_response_time{scenario:light_load}': ['p(95)<300'],
        'get_user_response_time{scenario:light_load}': ['p(95)<300'],
        'get_user_accounts_response_time{scenario:light_load}': ['p(95)<300'],
        'get_user_debit_cards_response_time{scenario:light_load}': ['p(95)<300'],
        'get_user_saved_accounts_response_time{scenario:light_load}': ['p(95)<300'],
        'get_user_banners_response_time{scenario:light_load}': ['p(95)<300'],
        'http_req_duration{scenario:light_load}': ['p(95)<300'],

        // Normal Load
        'login_response_time{scenario:normal_load}': ['p(95)<3000'],
        'get_user_response_time{scenario:normal_load}': ['p(95)<3000'],
        'get_user_accounts_response_time{scenario:normal_load}': ['p(95)<3000'],
        'get_user_debit_cards_response_time{scenario:normal_load}': ['p(95)<3000'],
        'get_user_saved_accounts_response_time{scenario:normal_load}': ['p(95)<3000'],
        'get_user_banners_response_time{scenario:normal_load}': ['p(95)<3000'],
        'http_req_duration{scenario:normal_load}': ['p(95)<3000'],

        // Heavy Load
        'login_response_time{scenario:heavy_load}': ['p(95)<8000'],
        'get_user_response_time{scenario:heavy_load}': ['p(95)<8000'],
        'get_user_accounts_response_time{scenario:heavy_load}': ['p(95)<8000'],
        'get_user_debit_cards_response_time{scenario:heavy_load}': ['p(95)<8000'],
        'get_user_saved_accounts_response_time{scenario:heavy_load}': ['p(95)<8000'],
        'get_user_banners_response_time{scenario:heavy_load}': ['p(95)<8000'],
        'http_req_duration{scenario:heavy_load}': ['p(95)<8000'],

        // Error rate tolerance per scenario
        'error_rate{scenario:light_load}': ['rate<0.01'],
        'error_rate{scenario:normal_load}': ['rate<0.03'],
        'error_rate{scenario:heavy_load}': ['rate<0.05'],

        // Transaction thresholds (complete auth + dashboard cycles)
        'transaction_counter{scenario:light_load}': ['rate>1'],
        'transaction_counter{scenario:normal_load}': ['rate>10'],
        'transaction_counter{scenario:heavy_load}': ['rate>30'],

        // HTTP requests per second (all requests combined)
        'http_reqs{scenario:light_load}': ['rate>1'],
        'http_reqs{scenario:normal_load}': ['rate>10'],
        'http_reqs{scenario:heavy_load}': ['rate>30']
    },
};

// Load test users from file
const TEST_USERS = new SharedArray('users', function () {
    try {
        const file = open('./users.txt');
        return file
            .split('\n')
            .filter(Boolean)
            .map((userId) => ({userId: userId.trim(), pin: '123456'}));
    } catch (err) {
        console.error('Failed to load users.txt:', err);
        // Fallback test users if file doesn't exist
        return [
            {userId: 'ffffd8dee1a111ef95a30242ac180002', pin: '123456'},
            {userId: 'ffff2e96e1a111ef95a30242ac180002', pin: '123456'},
            {userId: 'fffd93ece1a111ef95a30242ac180002', pin: '123456'}
        ];
    }
});

function getRandomUser() {
    return TEST_USERS[Math.floor(Math.random() * TEST_USERS.length)];
}

function getUser() {
    const user = getRandomUser();
    const payload = JSON.stringify({user_id: user.userId});
    const params = {
        headers: {'Content-Type': 'application/json'}
    };

    const start = Date.now();
    const res = http.post(`${BASE_URL}/api/v1/get-user-by-id`, payload, params);
    const duration = Date.now() - start;
    getUserResponseTime.add(duration);

    const ok = check(res, {
        'get_user status is 200': (r) => r.status === 200,
        'get_user returns user_info': (r) => {
            try {
                const data = JSON.parse(r.body);
                return data.data && data.data.user_info.name;
            } catch (_) {
                return false;
            }
        },
    });

    errorRate.add(!ok);
    if (!ok) {
        return null;
    }

    successfulGetUser.add(1);
    return user;
}

function login(user) {
    const payload = JSON.stringify({user_id: user.userId, pin: user.pin});
    const params = {
        headers: {'Content-Type': 'application/json'}
    };

    const start = Date.now();
    const res = http.post(`${BASE_URL}/api/v1/login`, payload, params);
    const duration = Date.now() - start;
    loginResponseTime.add(duration);

    const ok = check(res, {
        'login status is 200': (r) => r.status === 200,
        'login returns token': (r) => {
            try {
                const data = JSON.parse(r.body);
                return data.data && data.data.token;
            } catch (_) {
                return false;
            }
        },
    });

    errorRate.add(!ok);
    if (!ok) {
        return null;
    }

    successfulLogin.add(1);
    const data = JSON.parse(res.body).data;
    return {token: data.token};
}

function getUserAccounts(token) {
    if (!token?.token) return false;

    const params = {
        headers: {
            Authorization: `Bearer ${token.token}`,
        },
    };

    const start = Date.now();
    const res = http.get(`${BASE_URL}/api/v1/get-user-accounts`, params);
    const duration = Date.now() - start;
    getUserAccountsResponseTime.add(duration);

    const ok = check(res, {
        'get_user_accounts status is 200': (r) => r.status === 200,
        'get_user_accounts returns data': (r) => {
            try {
                const data = JSON.parse(r.body);
                return data.data !== undefined;
            } catch (_) {
                return false;
            }
        },
    });

    errorRate.add(!ok);
    if (!ok) {
        return false;
    } else {
        successfulGetUserAccounts.add(1);
        return true;
    }
}

function getUseDebitCards(token) {
    if (!token?.token) return false;

    const params = {
        headers: {
            Authorization: `Bearer ${token.token}`,
        },
    };

    const start = Date.now();
    const res = http.get(`${BASE_URL}/api/v1/get-user-debit-cards`, params);
    const duration = Date.now() - start;
    getUseDebitCardsResponseTime.add(duration);

    const ok = check(res, {
        'get_user_debit_cards status is 200': (r) => r.status === 200,
        'get_user_debit_cards returns data': (r) => {
            try {
                const data = JSON.parse(r.body);
                return data.data !== undefined;
            } catch (_) {
                return false;
            }
        },
    });

    errorRate.add(!ok);
    if (!ok) {
        return false;
    } else {
        successfulGetUserDebitCards.add(1);
        return true;
    }
}

function getUserSavedAccounts(token) {
    if (!token?.token) return false;

    const params = {
        headers: {
            Authorization: `Bearer ${token.token}`,
        },
    };

    const start = Date.now();
    const res = http.get(`${BASE_URL}/api/v1/get-user-saved-accounts`, params);
    const duration = Date.now() - start;
    getUserSavedAccountsResponseTime.add(duration);

    const ok = check(res, {
        'get_user_saved_accounts status is 200': (r) => r.status === 200,
        'get_user_saved_accounts returns data': (r) => {
            try {
                const data = JSON.parse(r.body);
                return data.data !== undefined;
            } catch (_) {
                return false;
            }
        },
    });

    errorRate.add(!ok);
    if (!ok) {
        return false;
    } else {
        successfulGetUserSavedAccounts.add(1);
        return true;
    }
}

function getUserBanners(token) {
    if (!token?.token) return false;

    const params = {
        headers: {
            Authorization: `Bearer ${token.token}`,
        },
    };

    const start = Date.now();
    const res = http.get(`${BASE_URL}/api/v1/get-user-banners`, params);
    const duration = Date.now() - start;
    getUserBannersResponseTime.add(duration);

    const ok = check(res, {
        'get_user_banners status is 200': (r) => r.status === 200,
        'get_user_banners returns data': (r) => {
            try {
                const data = JSON.parse(r.body);
                return data.data !== undefined;
            } catch (_) {
                return false;
            }
        },
    });

    errorRate.add(!ok);
    if (!ok) {
        return false;
    } else {
        successfulGetUserBanners.add(1);
        return true;
    }
}

export default function () {
    const user = getUser();
    const token = login(user);
    if (token) {
        const getUserAccountsSuccess = getUserAccounts(token);
        const getUseDebitCardsSuccess = getUseDebitCards(token);
        const getUserSavedAccountsSuccess = getUserSavedAccounts(token);
        const getUserBannersSuccess = getUserBanners(token);

        // Count as complete transaction only if all api calls succeed
        if (getUserAccountsSuccess && getUseDebitCardsSuccess && getUserSavedAccountsSuccess && getUserBannersSuccess) {
            transactionCounter.add(1);
        }
    }

    // Random sleep between 1-3 seconds to simulate realistic user behavior
    sleep(Math.random() * 2 + 1);
}

export function handleSummary(data) {
    // Calculate actual TPS from the summary data
    // const scenarios = ['light_load', 'normal_load', 'heavy_load'];
    const scenarios = ['light_load'];

    console.log('\n=== TPS Summary ===');
    scenarios.forEach(scenario => {
        const httpReqs = data.metrics.http_reqs?.values;
        const transactions = data.metrics.transaction_counter?.values;

        if (httpReqs && httpReqs.rate) {
            console.log(`${scenario} - HTTP Requests/sec: ${httpReqs.rate.toFixed(2)}`);
        }

        if (transactions && transactions.rate) {
            console.log(`${scenario} - Complete Transactions/sec: ${transactions.rate.toFixed(2)}`);
        }
    });

    return {
        'summary.html': htmlReport(data, {
            title: 'K6 Load Test Report - Auth & Dashboard Performance'
        }),
        stdout: textSummary(data, {indent: ' ', enableColors: true}),
    };
}