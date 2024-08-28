package main

import (
	"errors"
	"io"
	"net/http"
)

func getHTML(rawURL string) (string, error) {

	resp, err := http.Get(rawURL)

	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", errors.New(resp.Status)
	}

	if resp.Header.Get("Content-Type") == "text/html" {
		return "", errors.New("not text/html")
	}
	content, err := io.ReadAll(resp.Body)
	closingError := resp.Body.Close()
	if closingError != nil {
		return "", errors.New(closingError.Error())
	}
	if err != nil {
		return "", err
	}

	return string(content), nil
}
