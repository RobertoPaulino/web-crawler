package crawler

import (
	"fmt"
	"net/url"
	"sync"
)

// PageSEOData stores SEO-related information for a page
type PageSEOData struct {
	Title           string
	MetaDescription string
	H1Tags          []string
	StatusCode      int
	WordCount       int
	InternalLinks   int
	ExternalLinks   int
	ImagesCount     int
	HasCanonical    bool
	CanonicalURL    string
}

// Config represents the configuration and state for a web crawler
type Config struct {
	Pages              map[string]int
	SEOData            map[string]*PageSEOData
	BaseURL            *url.URL
	Mu                 *sync.Mutex
	ConcurrencyControl chan struct{}
	Wg                 *sync.WaitGroup
	MaxPages           int
}

// AddPageVisit adds a page to the visited pages map and returns whether it's the first visit
func (cfg *Config) AddPageVisit(normalizedURL string) (isFirst bool) {
	cfg.Mu.Lock()
	defer cfg.Mu.Unlock()

	if _, visited := cfg.Pages[normalizedURL]; visited {
		cfg.Pages[normalizedURL]++
		return false
	}

	cfg.Pages[normalizedURL] = 1
	return true
}

// StoreSEOData stores SEO data for a page
func (cfg *Config) StoreSEOData(normalizedURL string, data *PageSEOData) {
	cfg.Mu.Lock()
	defer cfg.Mu.Unlock()

	cfg.SEOData[normalizedURL] = data
}

// PagesLen returns the number of unique pages visited
func (cfg *Config) PagesLen() int {
	cfg.Mu.Lock()
	defer cfg.Mu.Unlock()
	return len(cfg.Pages)
}

// NewConfig creates a new crawler configuration
func NewConfig(rawBaseURL string, maxConcurrency int, maxPages int) (*Config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &Config{
		Pages:              make(map[string]int),
		SEOData:            make(map[string]*PageSEOData),
		BaseURL:            baseURL,
		Mu:                 &sync.Mutex{},
		ConcurrencyControl: make(chan struct{}, maxConcurrency),
		Wg:                 &sync.WaitGroup{},
		MaxPages:           maxPages,
	}, nil
}
