package scraper

import (
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Scraper struct {
	proxyManager *ProxyManager
}

func NewScraper(proxyManager *ProxyManager) *Scraper {
	return &Scraper{proxyManager: proxyManager}
}

func (s *Scraper) Scrape(url string) {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0"),
	)

	// Handle proxy
	proxy := s.proxyManager.GetRandomProxy()
	c.SetProxy(proxy)

	// Set rate limit to avoid bans
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		Delay:       2 * time.Second,
	})

	// Scrape logic
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		log.Println("Scraped link:", e.Attr("href"))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error:", err)
	})

	c.Visit(url)
	c.Wait()
}
