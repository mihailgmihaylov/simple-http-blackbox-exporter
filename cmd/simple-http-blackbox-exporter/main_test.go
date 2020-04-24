package main

import (
	"net/url"
	"testing"
)

func TestGetConfURLs(t *testing.T) {
	var c Config
	c.GetConf()
	for _, httpURL := range c.Urls {
		_, err := url.ParseRequestURI(httpURL)
		if err != nil {
			t.Errorf("Not a valid URL, got: %v, want: %v.", httpURL, "https://example.com/")
		}

		u, err := url.Parse(httpURL)
		if err != nil || u.Scheme == "" || u.Host == "" {
			t.Errorf("Incorrect url scheme or host: %v or %v.", u.Scheme, u.Host)
		}
	}
}

func TestGetURLTimeout(t *testing.T) {
	resp, _ := GetURL("https://httpstat.us/200?sleep=10000")
	if resp != 0 {
		t.Error("HTTP client get does not timeout for long queries.")
	}
}
