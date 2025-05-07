package crawler

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com")

	tests := []struct {
		name     string
		htmlBody string
		baseURL  *url.URL
		want     []string
		wantErr  bool
	}{
		{
			name:     "basic case",
			htmlBody: `<html><body><a href="/test">Test</a></body></html>`,
			baseURL:  baseURL,
			want:     []string{"https://example.com/test"},
			wantErr:  false,
		},
		{
			name:     "absolute URL",
			htmlBody: `<html><body><a href="https://google.com">Google</a></body></html>`,
			baseURL:  baseURL,
			want:     []string{"https://google.com"},
			wantErr:  false,
		},
		{
			name:     "multiple links",
			htmlBody: `<html><body><a href="/test1">Test1</a><a href="/test2">Test2</a></body></html>`,
			baseURL:  baseURL,
			want:     []string{"https://example.com/test1", "https://example.com/test2"},
			wantErr:  false,
		},
		{
			name:     "invalid HTML",
			htmlBody: `<html><body><a href="/test1">Test1</a><a href="</body></html>`,
			baseURL:  baseURL,
			want:     []string{"https://example.com/test1"}, // Only expect the valid URL
			wantErr:  false,                                 // We handle invalid URLs gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetURLsFromHTML(tt.htmlBody, tt.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetURLsFromHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// For the invalid HTML test, we only check that the valid URL is included
			if tt.name == "invalid HTML" {
				found := false
				for _, url := range got {
					if url == tt.want[0] {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetURLsFromHTML() = %v, expected to contain %v", got, tt.want[0])
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetURLsFromHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
