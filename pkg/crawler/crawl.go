package crawler

import (
	"fmt"
	"net/url"

	"github.com/RobertoPaulino/web-crawler/internal/utils"
)

// CrawlPage crawls a single page and schedules its links for crawling
func (cfg *Config) CrawlPage(rawCurrentURL string) {
	cfg.ConcurrencyControl <- struct{}{}
	defer func() {
		<-cfg.ConcurrencyControl
		cfg.Wg.Done()
	}()

	if cfg.PagesLen() >= cfg.MaxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
		return
	}

	// skip other websites
	if currentURL.Hostname() != cfg.BaseURL.Hostname() {
		return
	}

	normalizedURL, err := utils.NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - normalizedURL: %v\n", err)
	}

	isFirst := cfg.AddPageVisit(normalizedURL)
	if !isFirst {
		return
	}

	fmt.Printf("crawling %s\n", rawCurrentURL)

	resp, err := GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - getHTML: %v\n", err)
		return
	}

	// Extract and store SEO data
	seoData := ExtractSEOData(resp.Body, currentURL, cfg.BaseURL)
	seoData.StatusCode = resp.StatusCode
	cfg.StoreSEOData(normalizedURL, seoData)

	nextURLs, err := GetURLsFromHTML(resp.Body, cfg.BaseURL)
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v\n", err)
	}

	for _, nextURL := range nextURLs {
		cfg.Wg.Add(1)
		go cfg.CrawlPage(nextURL)
	}
}
