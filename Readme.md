# **Scalable Web Scraper with API Execution**

This project implements a scalable web scraper that supports API-based execution for scraping multiple URLs concurrently. It is designed to handle over **2 million requests per day**, incorporating proxy rotation, error handling, and concurrency.

---

## **Features**
1. **API-Based Scraping**: Trigger scraping tasks dynamically via a POST API.
2. **Concurrency**: Handles multiple scraping requests in parallel using Go routines.
3. **Proxy Rotation**: Integrates free proxy management to avoid IP bans.
4. **Rate Limiting**: Prevents overloading target servers and ensures compliance with their rate limits.
5. **Scalable Architecture**: Designed to distribute load across multiple scraper instances.
6. **Load Testing**: Includes a simulation script to validate high-load scenarios.
7. **Monitoring**: Integrates monitoring and scaling tools to ensure reliability.

---

## **Pre-Requisites**
1. **Software**:
   - [Go 1.20+](https://golang.org/)
   - [k6 Load Testing Tool](https://k6.io/docs/getting-started/installation/)
   - Docker and Docker Compose (for deployment and scalability testing)
   - Proxy list (configure free proxies for testing)

2. **Dependencies**:
   - Install project dependencies using `go mod tidy`.

3. **Environment Variables**:
   - Configure a `.env` file for:
     ```plaintext
     PROXY_POOL=http://proxy1,http://proxy2,http://proxy3
     ```

---

## **Steps to Run the Service**

### **1. Clone the Repository**
```bash
git clone https://github.com/satriyoaji/scrape-test.git
cd scrape-test
```

### **2. Install Dependencies**
```bash
go mod tidy
```

### **4. Adjust Env Vars**
```bash
cp .env.example .env #then adjust your own env vars
```

### **4. Run the API Service**
```bash
go run main.go
```

### **4. Trigger a Scraping Task**
```bash
curl -X POST http://localhost:8080/api/scrape \
-H "Content-Type: application/json" \
-d '{
  "urls": [
    "https://map.naver.com/",
    "https://shopping.naver.com/",
    "https://shoppinglive.naver.com/"
  ]
}'
```

## **System Architecture**
The scraper is built to handle massive scaling for over 2 million requests per day. The architecture is designed to be modular and resilient, employing proven tools like Redis and Kubernetes for queuing and scaling.

### **Architecture Design**
```plain
                  +----------------------+
                  |  Client Requests     |
                  +----------------------+
                            |
                            v
              +-------------------------------+
              |           API Gateway         |
              +-------------------------------+
                            |
                            v
       +----------------------------------------------+
       |         Load Balancer (NGINX/AWS ELB)        |
       +----------------------------------------------+
                            |
            +----------------------+------------------+
            |                      |                  |
   +------------------+   +------------------+   +------------------+
   | Scraper Service  |   | Scraper Service  |   | Scraper Service  |
   | Instance 1       |   | Instance 2       |   | Instance 3       |
   +------------------+   +------------------+   +------------------+
            |                      |                  |
            +------------------------------------------+
                            v
       +----------------------------------------------+
       |              Distributed Queue               |
       |        (Redis/RabbitMQ for Tasks)            |
       +----------------------------------------------+
                            |
                            v
       +----------------------------------------------+
       |          Centralized Storage (DB)           |
       +----------------------------------------------+

```

## **Testing Simulation**
### 1. **Load Testing with k6**
To simulate over 2 million requests, use the k6 load testing tool.
**Example load_test.js Script:**
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 1000, // Number of virtual users
  duration: '30m', // Test duration
};

export default function () {
  const url = 'http://localhost:8080/api/scrape';
  const payload = JSON.stringify({
    urls: [
       "https://map.naver.com/",
       "https://shopping.naver.com/",
       "https://shoppinglive.naver.com/"
    ],
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'is status 200': (r) => r.status === 200,
    'is valid JSON': (r) => r.json('results') !== null,
  });

  sleep(1);
}
```

**Run Load Test**:
```bash
k6 run load_test.js
```

### 2. **Monitor Test Result**
- Verify API responses. 
- Ensure the service can handle concurrent requests without crashing.

----

### Monitoring and Scaling Tools
1. Monitoring:
- PrometheusMonitor API metrics and performance. 
- Grafana: Visualize system health and alerts.

2. Scaling:
- Use Kubernetes Horizontal Pod Autoscaler (HPA) to scale scraper instances based on CPU and memory usage.

## Deployment
### Docker Compose Deployment
Create a `docker-compose.yml` file to deploy the scraper service along with a Redis instance for task queuing.
```yaml
version: '3.8'
services:
  scraper-service:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PROXY_POOL=${PROXY_POOL}
    depends_on:
      - redis
  redis:
    image: redis:6.2
    ports:
      - "6379:6379"
```
**Build and Deploy:**
```bash
docker-compose up --build
```
---
## **Future Enhancements**
1. CAPTCHA Solving:
- Integrate services like 2Captcha or Anti-Captcha.
2. Advanced Proxy Rotation:
- Use geo-targeted proxies for better IP distribution.
3. Distributed Tracing:
- Add tools like Jaeger for debugging large-scale systems.
4. Performance Optimization:
- Use caching to prevent redundant scraping of the same URLs.
