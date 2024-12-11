package scraper

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Scraper interface {
	Scrape(url string) error
}

type DynamicScraperWithChromedp struct {
	proxyManager *ProxyManager
}
type SimpleScraperWithColly struct {
	proxyManager *ProxyManager
}

func NewScraper(useChromedp bool, proxyManager *ProxyManager) Scraper {
	if useChromedp {
		return &DynamicScraperWithChromedp{proxyManager: proxyManager}
	}
	return &SimpleScraperWithColly{proxyManager: proxyManager}
}

func (s *SimpleScraperWithColly) Scrape(url string) error {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0"),
	)

	// Handle proxy
	proxy := s.proxyManager.GetRandomProxy()
	err := c.SetProxy(proxy)
	if err != nil {
		return err
	}

	// Set rate limit to avoid bans
	err = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		Delay:       2 * time.Second,
	})
	if err != nil {
		return err
	}

	// ScrapeWithColly logic
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		log.Println("Scraped link:", e.Attr("href"))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error:", err)
	})

	err = c.Visit(url)
	if err != nil {
		return err
	}
	c.Wait()

	return nil
}

func (s *DynamicScraperWithChromedp) Scrape(url string) error {
	//defer wg.Done()

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
		return err
	}

	log.Printf("Successfully scraped %s", url)
	// Log or store `res` in a database or file
	return nil
}
