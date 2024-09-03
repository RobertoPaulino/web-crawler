package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {

	htmlReader := strings.NewReader(htmlBody)
	doc, err := html.Parse(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("Error parsing HTML: %s", err)
	}

	var urls []string
	var traverseNodes func(*html.Node)
	traverseNodes = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					href, err := url.Parse(a.Val)
					if err != nil {
						fmt.Printf("Error parsing href: %s : %v", a.Val, err)
						continue
					}

					resolvedUrl := baseURL.ResolveReference(href)
					urls = append(urls, resolvedUrl.String())
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverseNodes(child)
		}
	}
	traverseNodes(doc)

	return urls, nil
}
