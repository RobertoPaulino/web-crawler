# Web Crawler

A concurrent web crawler written in Go that traverses websites, extracts links, and provides a report of all internal links found.

## Features

- Concurrent crawling with configurable concurrency limits
- Stays within the domain of the starting URL
- Configurable maximum number of pages to crawl
- Provides a detailed report with the count of internal links to each page
- Web interface with HTMX for interactive crawling
- Export results to CSV file
- **SEO Analysis**: Extracts and reports on key SEO metrics including:
  - Page titles and meta descriptions
  - H1 tags and content length
  - Internal and external link counts
  - Status codes and canonical URLs
  - Image counts and more

## Project Structure

```
web-crawler/
├── cmd/
│   ├── crawler/         # CLI application entry point
│   └── web/             # Web server entry point
├── pkg/
│   └── crawler/         # Core crawler functionality
├── internal/
│   └── utils/           # Internal utility functions
├── web/
│   ├── templates/       # HTML templates
│   └── static/          # Static assets (CSS, downloads)
├── Makefile             # Build and run commands
├── go.mod               # Go module definition
└── README.md            # Project documentation
```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/RobertoPaulino/web-crawler.git
   cd web-crawler
   ```

2. Build the project:
   ```bash
   make build
   ```

## Usage

### Command Line Interface

Run the crawler with the following command:

```bash
make run ARGS="<baseURL> <maxConcurrency> <maxPages>"
```

Example:
```bash
make run ARGS="https://example.com 10 100"
```

Or run directly:
```bash
./bin/crawler https://example.com 10 100"
```

#### Parameters

- `baseURL`: The starting URL for the crawler
- `maxConcurrency`: Maximum number of concurrent HTTP requests
- `maxPages`: Maximum number of pages to crawl

### Web Interface

Start the web server with:

```bash
make web
```

Then open your browser and navigate to:

```
http://localhost:8080
```

The web interface provides:
- A sleek, modern UI with gradient background
- Form to enter crawl parameters
- Interactive results display with SEO metrics
- Tabbed interface to view basic and detailed SEO data
- Option to export results as CSV with comprehensive SEO data

## SEO Analysis

The crawler collects the following SEO metrics for each page:

| Metric | Description |
|--------|-------------|
| Status Code | HTTP response code (200, 404, etc.) |
| Title | Page title from `<title>` tag |
| Meta Description | Content from meta description |
| Word Count | Total words on the page |
| Internal Links | Count of links pointing to the same domain |
| External Links | Count of links pointing to external domains |
| Images | Number of images on the page |
| H1 Tags | Number of H1 headings |
| Canonical URL | Self-referencing canonical URL if present |

These metrics are crucial for SEO analysis and can help identify:
- Missing title tags or meta descriptions
- Pages with duplicate content issues
- Pages with thin content
- Proper internal linking structure
- Crawlability issues

## Development

### Running Tests

```bash
make test
```

### Clean Build Artifacts

```bash
make clean
```

## License

This project is licensed under the terms found in the LICENSE file.