# Scraper Service

## Overview
This project is a scalable web scraping microservice designed to scrape websites like Shopee.tw or similar JavaScript-rendered websites. It uses two approaches:

1. **Colly**: Lightweight and fast, ideal for static or simple pages.
2. **Chromedp**: Headless browser for scraping JavaScript-heavy and dynamic pages.

The service supports the use of proxy managers, concurrent scraping, and rate-limiting to handle high-volume scraping requests without getting blocked or banned.

---

## Features
- **Dynamic Scraper Selection**: Automatically choose between `colly` or `chromedp` based on the website's complexity.
- **Proxy Management**: Rotates proxies to avoid IP bans and distribute requests.
- **Concurrency**: Leverages Go’s goroutines to scrape multiple URLs simultaneously.
- **Rate-Limiting**: Implements request throttling to prevent rate-limit bans.
- **Scalable Design**: Designed to handle over 2 million requests per day with support for tools like Redis, RabbitMQ, and Kubernetes.
- **Error Handling**: Logs errors and retries failed requests.

---

## System Architecture
The scraping service is designed as a microservice with scalability and fault tolerance in mind. Below is the high-level architecture:

### Diagram
```text
+-----------------------+    +-----------------------+
|   API Gateway        | → |   Load Balancer       |
+-----------------------+    +-----------------------+
            |                           |
  +-------------------+       +-------------------+
  | Scraper Instance 1|       | Scraper Instance 2|
  +-------------------+       +-------------------+
            |                           |
+----------------------+      +-----------------------+
| Proxy Manager (Redis)|      | Proxy Manager (Redis)|
+----------------------+      +-----------------------+
            |                           |
     +-----------------------+
     | Database/Log Storage  |
     +-----------------------+
```

### Components

1. **API Gateway**:
    - Exposes endpoints to trigger scraping tasks and retrieve results.
    - Routes traffic to scraper instances.

2. **Load Balancer**:
    - Distributes traffic evenly across multiple scraper instances.
    - Ensures horizontal scalability.

3. **Scraper Service**:
    - Core logic for scraping websites using `colly` or `chromedp`.
    - Handles concurrent scraping with goroutines.
    - Rotates proxies and implements rate limiting.

4. **Proxy Manager**:
    - Manages proxy pools for requests.
    - Uses Redis to cache and rotate proxies efficiently.

5. **Database/Log Storage**:
    - Stores scraped data, logs, and error reports for analysis.

---

## Setup
### Prerequisites
1. **Go (Golang)**: Install the latest version of Go.
2. **Redis**: For managing proxies and queues.
3. **Chromedp**: Ensure a headless Chrome browser is installed.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/scraper-service.git
   cd scraper-service
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run Redis:
   ```bash
   docker run -d -p 6379:6379 redis
   ```

4. Start the service:
   ```bash
   go run main.go
   ```

---

## Usage
### Configuration
Update the `.env` file with your configuration:
```env
REDIS_URL=localhost:6379
PROXY_POOL=your_proxy_pool
RATE_LIMIT=5   # Max requests per second per scraper
```

### Trigger Scraping
Modify the `urls` array in `main.go`:
```go
urls := []string{"https://shopee.tw", "https://shopee.tw/products", "https://shopee.tw/deals"}
```
Run the scraper:
```bash
go run main.go
```
Logs will show the scraped results or errors.

---

## Scalability
To handle more than 2 million requests per day:

1. **Horizontal Scaling**:
    - Deploy multiple scraper instances in a Kubernetes cluster.
    - Use a load balancer to distribute traffic.

2. **Task Queue**:
    - Use RabbitMQ or Redis as a queue manager to handle incoming scrape requests.
    - Allow tasks to be processed asynchronously by multiple scraper instances.

3. **Proxy Pooling**:
    - Integrate a proxy service with a large pool of IP addresses.
    - Rotate proxies dynamically to avoid bans.

4. **Dynamic Rate-Limiting**:
    - Adjust scraping speed based on website response to avoid triggering rate limits.

5. **Monitoring and Logging**:
    - Use tools like Prometheus and Grafana for monitoring metrics (e.g., request success rate, errors).
    - Store logs in a centralized system (e.g., ELK Stack).

---

## Testing with 2 Million Requests
1. Simulate requests using a task queue:
    - Push 2 million URLs into RabbitMQ or Redis queue.

2. Modify the scraper to fetch URLs from the queue:
```go
func fetchFromQueue(queue QueueManager) {
    for {
        url := queue.Pop()
        if url == "" {
            continue
        }
        scrape(url)
    }
}
```

3. Run multiple scraper instances to process the queue concurrently.
4. Monitor the progress and performance metrics.

---

## Future Enhancements
1. **Captcha Handling**:
    - Integrate third-party captcha-solving services like 2Captcha.

2. **Headless Browser Pooling**:
    - Manage a pool of headless browser instances to reduce startup overhead.

3. **Adaptive Scraping**:
    - Implement AI/ML to detect and adapt to anti-scraping mechanisms dynamically.

---

## Contributing
Contributions are welcome! Please create an issue or submit a pull request.

---

## License
This project is licensed under the MIT License. See the LICENSE file for details.

