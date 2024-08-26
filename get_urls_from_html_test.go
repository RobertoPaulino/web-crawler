package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		inputBody     string
		expected      []string
		errorContains string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<a href="/path/one">
							<span>Boot.dev</span>
						</a>
						<a href="https://other.com/path/one">
							<span>Boot.dev</span>
						</a>
					</body>
				</html>
				`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "relative URLs with query parameters",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
				<a href="/search?q=golang">
					<span>Search Golang</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://example.com/search?q=golang"},
		},
		{
			name:     "absolute URL with fragment",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
				<a href="https://example.com/path#section2">
					<span>Section 2</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://example.com/path#section2"},
		},
		{
			name:     "protocol-relative URLs",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
				<a href="//cdn.example.com/assets">
					<span>Assets</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://cdn.example.com/assets"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.errorContains == "" {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error containing '%v', got none.", i, tc.name, tc.errorContains)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected URLs %v, got URLs %v", i, tc.name, tc.expected, actual)
				return
			}
		})
	}
}
