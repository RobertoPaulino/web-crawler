package crawler

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPResponse contains the HTML body and additional metadata from the HTTP response
type HTTPResponse struct {
	Body       string
	StatusCode int
	Headers    http.Header
}

// GetHTML fetches the HTML content from a URL and returns it along with response metadata
func GetHTML(urlString string) (HTTPResponse, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(urlString)
	if err != nil {
		return HTTPResponse{StatusCode: 0}, fmt.Errorf("failed to get page: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HTTPResponse{StatusCode: resp.StatusCode, Headers: resp.Header}, fmt.Errorf("failed to read response body: %v", err)
	}

	return HTTPResponse{
		Body:       string(body),
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}
