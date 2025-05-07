package utils

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "basic case",
			url:     "https://example.com/path",
			want:    "example.com/path",
			wantErr: false,
		},
		{
			name:    "with trailing slash",
			url:     "https://example.com/path/",
			want:    "example.com/path/",
			wantErr: false,
		},
		{
			name:    "with query params",
			url:     "https://example.com/path?query=value",
			want:    "example.com/path",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			url:     "://invalid",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
