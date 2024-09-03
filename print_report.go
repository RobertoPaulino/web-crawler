package main

import "fmt"

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("  REPORT for %v \n", baseURL)
	fmt.Println("=============================")

	report := mapSort(pages)
	for _, str := range report {
		fmt.Print(str)
	}
}
