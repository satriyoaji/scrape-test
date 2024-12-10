package tests

import (
	"context"
	"log"
	"testing"

	"webscraper/scraper"
)

func TestLoad(t *testing.T) {
	//config.LoadConfig()
	//fmt.Println("env: ", config.GetEnv("REDIS_URL"))
	queueManager := scraper.NewQueueManager("localhost:6379")
	ctx := context.Background()

	// Push 2M URLs
	for i := 0; i < 2000000; i++ {
		queueManager.Push(ctx, "scrape_queue", "https://shopee.tw/product")
	}

	log.Println("2M URLs added to queue for load testing.")
}
