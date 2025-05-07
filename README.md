# Web Crawler (Vercel Deployment Branch)

This branch is optimized for deployment on [Vercel](https://vercel.com). It contains only the files and configuration needed to run the web crawler as a serverless Go function with a modern web interface.

## Features

- Deployable instantly to Vercel as a serverless Go function
- Interactive web interface for crawling and SEO research
- Downloadable CSV with detailed SEO data for each crawled page
- Clean, modern UI

## Usage

1. **Deploy to Vercel**
   - Clone this repository or connect it to your Vercel account.
   - Deploy directly from the Vercel dashboard or using the CLI:
     ```bash
     vercel --prod
     ```
   - No build steps or shell scripts are required.

2. **Using the Web Crawler**
   - Visit your deployed Vercel URL.
   - Enter the URL you want to crawl, set concurrency and max pages if desired, and start crawling.
   - Results will appear in a clean table.
   - For full SEO data, click **Export SEO Data as CSV** to download a spreadsheet with all metrics.

## Project Structure

```
web-crawler/
├── cmd/
│   └── web/             # Vercel Go serverless function (entry point)
├── pkg/
│   └── crawler/         # Core crawler functionality
├── internal/
│   └── utils/           # Internal utility functions
├── static/
│   └── css/             # Static CSS assets
├── vercel.json          # Vercel deployment configuration
├── go.mod               # Go module definition
└── README.md            # This file
```

## Notes

- This branch is for deployment only. For development, use the main branch.
- All shell scripts and Makefiles have been removed for a clean, serverless deployment.
- The crawler runs in a serverless environment, so each crawl is stateless and limited by Vercel's function execution time.
- For advanced SEO analysis, always use the CSV export feature.

## License

This project is licensed under the terms found in the LICENSE file.