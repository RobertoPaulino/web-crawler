package main

import (
	"errors"
	"net/url"
)

func normalizeURL(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	print(parsedURL)
	if u == "" {
		return "", errors.New("normalizeURL: URL is empty")
	}

	if err != nil || (parsedURL.Host == "" && parsedURL.Path == "") {
		return "", errors.New("normalizeURL: URL is invalid")
	}

	return parsedURL.Host + parsedURL.Path, nil
}
