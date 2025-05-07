package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/RobertoPaulino/web-crawler/pkg/crawler"
)

// PageResult represents a single page with link count and SEO data
type PageResult struct {
	URL             string
	Count           int
	Title           string
	StatusCode      int
	WordCount       int
	InternalLinks   int
	ExternalLinks   int
	ImagesCount     int
	HasH1           bool
	H1Count         int
	HasMeta         bool
	HasCanonical    bool
	CanonicalURL    string
	MetaDescription string
}

// ResultsData holds the data to be passed to the results template
type ResultsData struct {
	BaseURL     string
	Concurrency int
	MaxPages    int
	PageCount   int
	Results     []PageResult
	Error       string
}

// DownloadData holds the data for download template
type DownloadData struct {
	URL string
}

// Handler is the main HTTP handler for the application
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set up templates
	templates := template.Must(template.ParseGlob("web/templates/*.html"))

	// Serve static files
	if r.URL.Path == "/static/" {
		fs := http.FileServer(http.Dir("web/static"))
		fs.ServeHTTP(w, r)
		return
	}

	// Home page
	if r.URL.Path == "/" {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Handle crawl
	if r.URL.Path == "/crawl" {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		url := r.FormValue("url")
		concurrencyStr := r.FormValue("concurrency")
		pagesStr := r.FormValue("pages")

		concurrency, err := strconv.Atoi(concurrencyStr)
		if err != nil || concurrency < 1 {
			concurrency = 5 // Default
		}

		maxPages, err := strconv.Atoi(pagesStr)
		if err != nil || maxPages < 1 {
			maxPages = 50 // Default
		}

		// Configure and run crawler
		cfg, err := crawler.NewConfig(url, concurrency, maxPages)
		if err != nil {
			renderResults(w, templates, ResultsData{
				Error: fmt.Sprintf("Failed to configure crawler: %v", err),
			})
			return
		}

		// Crawl
		cfg.Wg.Add(1)
		go cfg.CrawlPage(url)
		cfg.Wg.Wait()

		// Prepare results
		results := prepareResults(cfg)

		renderResults(w, templates, ResultsData{
			BaseURL:     url,
			Concurrency: concurrency,
			MaxPages:    maxPages,
			PageCount:   len(cfg.Pages),
			Results:     results,
		})
		return
	}

	// Export CSV
	if r.URL.Path == "/export-csv" {
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
			http.Error(w, fmt.Sprintf("Failed to configure crawler: %v", err), http.StatusInternalServerError)
			return
		}

		// Crawl
		cfg.Wg.Add(1)
		go cfg.CrawlPage(url)
		cfg.Wg.Wait()

		// Create CSV file
		timestamp := time.Now().Format("20060102-150405")
		filename := fmt.Sprintf("seo_crawl_%s.csv", timestamp)
		csvPath := filepath.Join("web/static/downloads", filename)

		// Ensure directory exists
		os.MkdirAll(filepath.Dir(csvPath), 0755)

		file, err := os.Create(csvPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create CSV file: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
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
			http.Error(w, fmt.Sprintf("Failed to write CSV header: %v", err), http.StatusInternalServerError)
			return
		}

		// Prepare results
		results := prepareResults(cfg)

		// Write data
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
				http.Error(w, fmt.Sprintf("Failed to write CSV data: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Return download link
		downloadURL := fmt.Sprintf("/static/downloads/%s", filename)
		err = templates.ExecuteTemplate(w, "download.html", DownloadData{URL: downloadURL})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Handle 404
	http.NotFound(w, r)
}

func prepareResults(cfg *crawler.Config) []PageResult {
	results := make([]PageResult, 0, len(cfg.Pages))

	for url, count := range cfg.Pages {
		result := PageResult{
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

func renderResults(w http.ResponseWriter, tmpl *template.Template, data ResultsData) {
	err := tmpl.ExecuteTemplate(w, "results.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
