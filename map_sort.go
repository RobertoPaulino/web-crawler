package main

import (
	"fmt"
	"sort"
)

func mapSort(pages map[string]int) []string {

	keys := make([]string, 0, len(pages))

	res := make([]string, len(pages))

	for key := range pages {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return pages[keys[i]] > pages[keys[j]]
	})

	for _, k := range keys {
		if pages[k] > 0 {
			res = append(res, fmt.Sprintf("Found %v internal links to %v \n", pages[k], k))
		}
	}

	return res
}
