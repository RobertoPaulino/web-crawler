package utils

import (
	"fmt"
	"sort"
)

func MapSort(pages map[string]int) []string {
	keys := make([]string, 0, len(pages))

	for key := range pages {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return pages[keys[i]] > pages[keys[j]]
	})

	var res []string
	for _, k := range keys {
		if pages[k] > 0 {
			res = append(res, fmt.Sprintf("Found %v internal links to %v \n", pages[k], k))
		}
	}

	return res
}
