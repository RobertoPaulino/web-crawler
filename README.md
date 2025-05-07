# Web Crawler

A concurrent web crawler written in Go that traverses websites, extracts links, and provides a report of all internal links found. Now with serverless deployment support for Vercel!

## Features

- Concurrent crawling with configurable concurrency limits
- Stays within the domain of the starting URL
- Configurable maximum number of pages to crawl
- Provides a detailed report with the count of internal links to each page
- Modern web interface with real-time updates
- Export results to CSV file
- **SEO Analysis**: Extracts and reports on key SEO metrics including:
  - Page titles and meta descriptions
  - H1 tags and content length
  - Internal and external link counts
  - Status codes and canonical URLs
  - Image counts and more
- **Serverless Deployment**: Ready to deploy on Vercel

## Project Structure

```
web-crawler/
├── api/              # Serverless API functions
│   └── index.go      # Main API handler
├── pkg/
│   └── crawler/      # Core crawler functionality
├── web/
│   └── templates/    # HTML templates
├── vercel.json       # Vercel configuration
├── go.mod           # Go module definition
└── README.md        # Project documentation
```

## Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/RobertoPaulino/web-crawler.git
   cd web-crawler
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the development server:
   ```bash
   go run api/index.go
   ```

4. Open your browser and navigate to:
   ```
   http://localhost:3000
   ```

## Deploying to Vercel

1. Install the Vercel CLI:
   ```bash
   npm i -g vercel
   ```

2. Login to Vercel:
   ```bash
   vercel login
   ```

3. Deploy the project:
   ```bash
   vercel
   ```

4. For production deployment:
   ```bash
   vercel --prod
   ```

The project will be automatically built and deployed to Vercel's serverless platform. The Go functions will be compiled and optimized for serverless execution.

## API Endpoints

### POST /api/crawl
Crawls a website and returns the results.

Request body:
```json
{
  "url": "https://example.com",
  "concurrency": 5,
  "maxPages": 50
}
```

Response:
```json
{
  "baseURL": "https://example.com",
  "concurrency": 5,
  "maxPages": 50,
  "pageCount": 10,
  "results": [
    {
      "url": "https://example.com",
      "count": 5,
      "title": "Example Domain",
      "statusCode": 200,
      "wordCount": 100,
      "internalLinks": 4,
      "externalLinks": 1,
      "imagesCount": 2,
      "hasH1": true,
      "h1Count": 1,
      "hasMeta": true,
      "hasCanonical": true,
      "canonicalURL": "https://example.com",
      "metaDescription": "Example domain for testing"
    }
  ]
}
```

### GET /api/export-csv
Exports crawl results as a CSV file.

Query parameters:
- `url`: The base URL to crawl
- `concurrency`: Maximum number of concurrent requests (default: 5)
- `pages`: Maximum number of pages to crawl (default: 50)

## License

This project is licensed under the terms found in the LICENSE file.