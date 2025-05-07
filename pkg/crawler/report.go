package crawler

import (
	"fmt"

	"github.com/RobertoPaulino/web-crawler/internal/utils"
)

// PrintReport prints a formatted report of the crawl results
func PrintReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("  REPORT for %v \n", baseURL)
	fmt.Println("=============================")

	report := utils.MapSort(pages)
	for _, str := range report {
		fmt.Print(str)
	}
}
