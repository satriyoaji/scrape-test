package main

import (
	"context"
	"fmt"
	"log"
	_ "strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"webscraper/config"
	"webscraper/scraper"
)

func scrape(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", &res),
	)
	if err != nil {
		log.Printf("Error scraping %s: %v", url, err)
		return
	}

	log.Printf("Successfully scraped %s", url)
	// Log or store `res` in a database or file
}

func main() {
	urls := []string{"https://shopee.tw", "https://shopee.tw/products", "https://shopee.tw/deals"}
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go scrape(url, &wg)
	}

	wg.Wait()
	log.Println("Scraping completed.")
}

func main2() {
	config.LoadConfig()

	// Initialize managers
	proxyManager := scraper.NewProxyManager(config.GetEnv("PROXY_POOL"))
	queueManager := scraper.NewQueueManager(config.GetEnv("REDIS_URL"))
	fmt.Println("Env: ", proxyManager, queueManager)
	ctx := context.Background()

	// Push URLs to queue (example)
	urls := []string{"https://shopee.tw", "https://shopee.tw/deals", "https://shopee.tw/products"}
	for _, url := range urls {
		queueManager.Push(ctx, "scrape_queue", url)
	}

	// Start workers
	numWorkers := 100
	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scraperService := scraper.NewScraper(proxyManager)
			for {
				url := queueManager.Pop(ctx, "scrape_queue")
				if url == "" {
					break
				}
				scraperService.Scrape(url)
			}
		}()
	}

	wg.Wait()
	log.Println("Scraping completed!")
}

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
