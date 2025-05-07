package api

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RobertoPaulino/web-crawler/pkg/crawler"
)

type CrawlRequest struct {
	URL         string `json:"url"`
	Concurrency int    `json:"concurrency"`
	MaxPages    int    `json:"maxPages"`
}

type CrawlResponse struct {
	BaseURL     string               `json:"baseURL"`
	Concurrency int                  `json:"concurrency"`
	MaxPages    int                  `json:"maxPages"`
	PageCount   int                  `json:"pageCount"`
	Results     []crawler.PageResult `json:"results"`
	Error       string               `json:"error,omitempty"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle different routes
	switch r.URL.Path {
	case "/":
		handleIndex(w, r)
	case "/api/crawl":
		handleCrawl(w, r)
	case "/api/export-csv":
		handleExportCSV(w, r)
	default:
		http.NotFound(w, r)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}

func handleCrawl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CrawlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set defaults if not provided
	if req.Concurrency < 1 {
		req.Concurrency = 5
	}
	if req.MaxPages < 1 {
		req.MaxPages = 50
	}

	// Configure and run crawler
	cfg, err := crawler.NewConfig(req.URL, req.Concurrency, req.MaxPages)
	if err != nil {
		json.NewEncoder(w).Encode(CrawlResponse{
			Error: "Failed to configure crawler: " + err.Error(),
		})
		return
	}

	// Crawl
	cfg.Wg.Add(1)
	go cfg.CrawlPage(req.URL)
	cfg.Wg.Wait()

	// Prepare results
	results := prepareResults(cfg)

	// Send response
	json.NewEncoder(w).Encode(CrawlResponse{
		BaseURL:     req.URL,
		Concurrency: req.Concurrency,
		MaxPages:    req.MaxPages,
		PageCount:   len(cfg.Pages),
		Results:     results,
	})
}

func handleExportCSV(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	url := r.URL.Query().Get("url")
	concurrencyStr := r.URL.Query().Get("concurrency")
	pagesStr := r.URL.Query().Get("pages")

	concurrency, err := strconv.Atoi(concurrencyStr)
	if err != nil || concurrency < 1 {
		concurrency = 5
	}

	maxPages, err := strconv.Atoi(pagesStr)
	if err != nil || maxPages < 1 {
		maxPages = 50
	}

	// Configure and run crawler
	cfg, err := crawler.NewConfig(url, concurrency, maxPages)
	if err != nil {
		http.Error(w, "Failed to configure crawler: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Crawl
	cfg.Wg.Add(1)
	go cfg.CrawlPage(url)
	cfg.Wg.Wait()

	// Set headers for CSV download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=seo_crawl.csv")

	// Write CSV
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	headers := []string{
		"URL",
		"Link Count",
		"Status Code",
		"Title",
		"Meta Description",
		"Word Count",
		"Internal Links",
		"External Links",
		"Images Count",
		"Has H1",
		"H1 Count",
		"Has Meta Description",
		"Has Canonical",
		"Canonical URL",
	}

	if err := writer.Write(headers); err != nil {
		http.Error(w, "Failed to write CSV header", http.StatusInternalServerError)
		return
	}

	// Write data
	results := prepareResults(cfg)
	for _, result := range results {
		row := []string{
			result.URL,
			strconv.Itoa(result.Count),
			strconv.Itoa(result.StatusCode),
			result.Title,
			result.MetaDescription,
			strconv.Itoa(result.WordCount),
			strconv.Itoa(result.InternalLinks),
			strconv.Itoa(result.ExternalLinks),
			strconv.Itoa(result.ImagesCount),
			strconv.FormatBool(result.HasH1),
			strconv.Itoa(result.H1Count),
			strconv.FormatBool(result.HasMeta),
			strconv.FormatBool(result.HasCanonical),
			result.CanonicalURL,
		}

		if err := writer.Write(row); err != nil {
			http.Error(w, "Failed to write CSV data", http.StatusInternalServerError)
			return
		}
	}
}

func prepareResults(cfg *crawler.Config) []crawler.PageResult {
	results := make([]crawler.PageResult, 0, len(cfg.Pages))

	for url, count := range cfg.Pages {
		result := crawler.PageResult{
			URL:   url,
			Count: count,
		}

		// Add SEO data if available
		if seoData, ok := cfg.SEOData[url]; ok {
			result.Title = seoData.Title
			result.StatusCode = seoData.StatusCode
			result.WordCount = seoData.WordCount
			result.InternalLinks = seoData.InternalLinks
			result.ExternalLinks = seoData.ExternalLinks
			result.ImagesCount = seoData.ImagesCount
			result.HasH1 = len(seoData.H1Tags) > 0
			result.H1Count = len(seoData.H1Tags)
			result.HasMeta = seoData.MetaDescription != ""
			result.MetaDescription = seoData.MetaDescription
			result.HasCanonical = seoData.HasCanonical
			result.CanonicalURL = seoData.CanonicalURL
		}

		results = append(results, result)
	}

	return results
}
