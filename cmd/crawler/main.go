package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/RobertoPaulino/web-crawler/pkg/crawler"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("not enough arguments provided")
		fmt.Println("usage: crawler <baseURL> <maxConcurrency> <maxPages>")
		return
	}
	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		return
	}
	rawBaseURL := os.Args[1]
	maxConcurrencyString := os.Args[2]
	maxPagesString := os.Args[3]

	maxConcurrency, err := strconv.Atoi(maxConcurrencyString)
	if err != nil {
		fmt.Printf("Error - maxConcurrency: %v\n", err)
		return
	}
	maxPages, err := strconv.Atoi(maxPagesString)
	if err != nil {
		fmt.Printf("Error - maxPages: %v\n", err)
		return
	}

	cfg, err := crawler.NewConfig(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error - configure: %v\n", err)
		return
	}

	fmt.Printf("starting crawl of: %s...\n", rawBaseURL)

	cfg.Wg.Add(1)
	go cfg.CrawlPage(rawBaseURL)
	cfg.Wg.Wait()

	crawler.PrintReport(cfg.Pages, rawBaseURL)
}
