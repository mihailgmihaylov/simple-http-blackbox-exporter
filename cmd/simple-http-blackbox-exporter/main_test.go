package main

import (
	"net/url"
	"strconv"
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

func TestGetURLValues(t *testing.T) {
	var codes = []int{200, 201, 202, 203, 204, 205, 206, 301, 302, 303, 304, 305, 306, 307, 308, 400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 500, 501, 502, 503, 504, 505, 506, 507}
	for _, code := range codes {
		url := "https://httpstat.us/" + strconv.Itoa(code)
		resp, time := GetURL(url)
		if resp == 0 && time != 0 {
			t.Errorf(
				"Incorrect response time - it should be 0 when response failed but it is: %v for URL: %v",
				time, url,
			)
		}

		if resp != 0 && resp != 1 {
			t.Errorf("Incorrect response code received: %v but it should be 0 or 1.", time)
		}

		if time < 0 {
			t.Errorf("Negative response time received: %v but it should be a positive value.", time)
		}
	}
}

func TestGetURLTimeout(t *testing.T) {
	resp, _ := GetURL("https://httpstat.us/200?sleep=10000")
	if resp != 0 {
		t.Error("HTTP client get does not timeout for long queries.")
	}
}
