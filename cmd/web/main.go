package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"net/http"
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

var (
	indexTemplate    = template.Must(template.New("index").Parse(indexHTML))
	resultsTemplate  = template.Must(template.New("results").Parse(resultsHTML))
	downloadTemplate = template.Must(template.New("download").Parse(downloadHTML))
)

// Handler is the main HTTP handler for the application
func Handler(w http.ResponseWriter, r *http.Request) {
	// Home page
	if r.URL.Path == "/" {
		err := indexTemplate.Execute(w, nil)
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
			renderResults(w, ResultsData{
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

		renderResults(w, ResultsData{
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

		// Create CSV in memory
		var buf bytes.Buffer
		writer := csv.NewWriter(&buf)

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

		writer.Flush()
		if err := writer.Error(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write CSV: %v", err), http.StatusInternalServerError)
			return
		}

		// Set headers for CSV download
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=seo_crawl_%s.csv", time.Now().Format("20060102-150405")))

		// Write CSV to response
		if _, err := io.Copy(w, &buf); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send CSV: %v", err), http.StatusInternalServerError)
			return
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

func renderResults(w http.ResponseWriter, data ResultsData) {
	err := resultsTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Embedded HTML templates
const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Web Crawler</title>
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
    <div class="container">
        <h1>Web Crawler</h1>
        <div class="crawler-form">
            <form hx-post="/crawl" hx-target="#results" hx-indicator=".spinner-container">
                <div class="form-group">
                    <input type="url" name="url" placeholder="Enter URL to crawl" required>
                </div>
                <div class="form-group controls">
                    <div>
                        <input type="number" name="concurrency" placeholder="Concurrency (default: 5)" min="1" max="20">
                    </div>
                    <div>
                        <input type="number" name="pages" placeholder="Max Pages (default: 50)" min="1" max="1000">
                    </div>
                </div>
                <button type="submit">Start Crawling</button>
            </form>
        </div>
        <div class="spinner-container htmx-indicator">
            <div class="spinner"></div>
            <p>Crawling in progress...</p>
        </div>
        <div id="results"></div>
    </div>
</body>
</html>`

const resultsHTML = `{{if .Error}}
    <div class="error">
        <p>Error: {{.Error}}</p>
    </div>
{{else}}
    <div class="crawl-info">
        <h2>Crawl Results for {{.BaseURL}}</h2>
        <p class="page-count">Found {{.PageCount}} unique pages</p>
    </div>
    
    {{if eq .PageCount 0}}
        <div class="empty-results">
            <p>No pages found.</p>
        </div>
    {{else}}
        <div class="results-tabs">
            <button class="tab-btn active" data-tab="basic">Basic</button>
            <button class="tab-btn" data-tab="seo">SEO Details</button>
        </div>

        <div id="basic-tab" class="tab-content active">
            <table>
                <thead>
                    <tr>
                        <th>URL</th>
                        <th>Links</th>
                        <th>Status</th>
                        <th>Title</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Results}}
                    <tr>
                        <td>{{.URL}}</td>
                        <td>{{.Count}}</td>
                        <td>{{.StatusCode}}</td>
                        <td>{{.Title}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>

        <div id="seo-tab" class="tab-content">
            <table>
                <thead>
                    <tr>
                        <th>URL</th>
                        <th>Words</th>
                        <th>Int. Links</th>
                        <th>Ext. Links</th>
                        <th>Images</th>
                        <th>H1</th>
                        <th>Meta</th>
                        <th>Canonical</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Results}}
                    <tr>
                        <td>{{.URL}}</td>
                        <td>{{.WordCount}}</td>
                        <td>{{.InternalLinks}}</td>
                        <td>{{.ExternalLinks}}</td>
                        <td>{{.ImagesCount}}</td>
                        <td class="{{if .HasH1}}good{{else}}bad{{end}}">{{.H1Count}}</td>
                        <td class="{{if .HasMeta}}good{{else}}bad{{end}}">{{if .HasMeta}}Yes{{else}}No{{end}}</td>
                        <td class="{{if .HasCanonical}}good{{else}}bad{{end}}">
                            {{if .HasCanonical}}
                                <span title="{{.CanonicalURL}}">Yes</span>
                            {{else}}
                                No
                            {{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        
        <div class="actions">
            <a href="/export-csv?url={{.BaseURL}}&concurrency={{.Concurrency}}&pages={{.MaxPages}}" class="button">
                Export SEO Data as CSV
            </a>
        </div>
    {{end}}

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            const tabBtns = document.querySelectorAll('.tab-btn');
            const tabContents = document.querySelectorAll('.tab-content');
            
            tabBtns.forEach(btn => {
                btn.addEventListener('click', () => {
                    const tabId = btn.getAttribute('data-tab');
                    
                    // Hide all content
                    tabContents.forEach(content => {
                        content.classList.remove('active');
                    });
                    
                    // Remove active class from all buttons
                    tabBtns.forEach(btn => {
                        btn.classList.remove('active');
                    });
                    
                    // Show selected content and mark button as active
                    document.getElementById(tabId + '-tab').classList.add('active');
                    btn.classList.add('active');
                });
            });
        });
    </script>
{{end}}`

const downloadHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Download Complete</title>
    <link rel="stylesheet" href="/static/css/styles.css">
</head>
<body>
    <div class="container">
        <h1>Download Complete</h1>
        <div class="download-info">
            <p>Your CSV file has been generated and should start downloading automatically.</p>
            <p>If the download doesn't start automatically, <a href="{{.URL}}">click here</a>.</p>
        </div>
    </div>
</body>
</html>`
