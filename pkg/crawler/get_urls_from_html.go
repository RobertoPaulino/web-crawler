package crawler

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// GetURLsFromHTML extracts all URLs from the HTML body that are within the same domain
func GetURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	var urls []string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					if href == "" || strings.HasPrefix(href, "#") {
						continue
					}

					u, err := url.Parse(href)
					if err != nil {
						continue
					}

					if u.Host == "" || u.Host == baseURL.Host {
						urls = append(urls, baseURL.ResolveReference(u).String())
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	return urls, nil
}
