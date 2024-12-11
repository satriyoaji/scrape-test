package main

import (
	"encoding/json"
	"log"
	"net/http"
	_ "strings"
	"sync"
	"webscraper/config"
	"webscraper/scraper"
)

type ScrapeRequest struct {
	URLs []string `json:"urls"`
}

type ScrapeResponse struct {
	Results map[string]string `json:"results"`
}

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req ScrapeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.URLs) == 0 {
		http.Error(w, "Invalid JSON payload or empty URLs", http.StatusBadRequest)
		return
	}

	results := make(map[string]string)
	var wg sync.WaitGroup

	proxyManager := scraper.NewProxyManager(config.GetEnv("PROXY_POOL"))

	for _, url := range req.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Use chromedp for dynamic scraping
			s := scraper.NewScraper(true, proxyManager)
			err := s.Scrape(url)
			if err != nil {
				results[url] = "Error: " + err.Error()
			} else {
				results[url] = "Successfully scraped"
			}
		}(url)
	}

	wg.Wait()

	resp := ScrapeResponse{Results: results}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	config.LoadConfig()

	http.HandleFunc("/api/scrape", scrapeHandler)

	log.Printf("Server running on port %s, \n", config.GetEnv("API_PORT"))
	err := http.ListenAndServe(config.GetEnv("API_PORT"), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

//func main2() {
//	config.LoadConfig()
//
//	// Initialize managers
//	proxyManager := scraper.NewProxyManager(config.GetEnv("PROXY_POOL"))
//	queueManager := scraper.NewQueueManager(config.GetEnv("REDIS_URL"))
//	fmt.Println("Env: ", proxyManager, queueManager)
//	ctx := context.Background()
//
//	// Push URLs to queue (example)
//	urls := []string{"https://shopee.tw", "https://shopee.tw/deals", "https://shopee.tw/products"}
//	for _, url := range urls {
//		queueManager.Push(ctx, "scrape_queue", url)
//	}
//
//	// Start workers
//	numWorkers := 100
//	wg := &sync.WaitGroup{}
//
//	for i := 0; i < numWorkers; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			scraperService := scraper.NewScraper(proxyManager)
//			for {
//				url := queueManager.Pop(ctx, "scrape_queue")
//				if url == "" {
//					break
//				}
//				scraperService.ScrapeWithColly(url)
//			}
//		}()
//	}
//
//	wg.Wait()
//	log.Println("Scraping completed!")
//}

//func main3() {
//	// Setup proxy
//	proxy := goproxy.NewProxyHttpServer()
//	proxy.Verbose = true
//	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}
//
//	// Create a collector
//	c := colly.NewCollector(
//		colly.Async(true), // Enable asynchronous scraping
//	)
//
//	c.WithTransport(httpClient.Transport)
//
//	// Handle the response
//	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
//		link := e.Attr("href")
//		fmt.Println("Found link:", link)
//	})
//
//	// Handle errors
//	c.OnError(func(r *colly.Response, err error) {
//		log.Println("Error:", err)
//	})
//
//	// Visit URLs
//	urls := []string{
//		"https://shopee.tw", // Add other URLs
//	}
//	for _, url := range urls {
//		err := c.Visit(url)
//		if err != nil {
//			log.Fatal(err)
//			return
//		}
//	}
//}
