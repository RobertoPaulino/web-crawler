package crawler

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ExtractSEOData extracts SEO-related information from HTML content
func ExtractSEOData(htmlContent string, pageURL *url.URL, baseURL *url.URL) *PageSEOData {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return &PageSEOData{}
	}

	data := &PageSEOData{
		InternalLinks: 0,
		ExternalLinks: 0,
		H1Tags:        []string{},
	}

	// Extract title, meta description, h1 tags, canonical, etc.
	extractNodes(doc, data, pageURL, baseURL)

	// Count words in the text content
	data.WordCount = countWords(extractTextContent(doc))

	return data
}

// extractNodes recursively traverses the HTML document and extracts SEO data
func extractNodes(n *html.Node, data *PageSEOData, pageURL *url.URL, baseURL *url.URL) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				data.Title = extractTextContent(n)
			}
		case "meta":
			// Extract meta description
			var name, content string
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == "description" {
					name = attr.Val
				}
				if attr.Key == "content" {
					content = attr.Val
				}
			}
			if name == "description" && content != "" {
				data.MetaDescription = content
			}
		case "h1":
			h1Text := extractTextContent(n)
			if h1Text != "" {
				data.H1Tags = append(data.H1Tags, h1Text)
			}
		case "link":
			// Check for canonical URL
			var rel, href string
			for _, attr := range n.Attr {
				if attr.Key == "rel" && attr.Val == "canonical" {
					rel = attr.Val
				}
				if attr.Key == "href" {
					href = attr.Val
				}
			}
			if rel == "canonical" && href != "" {
				data.HasCanonical = true
				data.CanonicalURL = href
			}
		case "a":
			// Count internal vs external links
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					if href == "" || strings.HasPrefix(href, "#") {
						continue
					}

					linkURL, err := url.Parse(href)
					if err != nil {
						continue
					}

					if linkURL.Host == "" || linkURL.Host == baseURL.Host {
						data.InternalLinks++
					} else {
						data.ExternalLinks++
					}
				}
			}
		case "img":
			data.ImagesCount++
		}
	}

	// Recursively process child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractNodes(c, data, pageURL, baseURL)
	}
}

// extractTextContent extracts all text content from an HTML node
func extractTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += extractTextContent(c)
	}
	return strings.TrimSpace(result)
}

// countWords counts the number of words in a text
func countWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}
