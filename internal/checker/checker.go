package checker

import (
	"time"
	"net/http"
)

type Result struct {
	URL        string
	StatusCode int
	OK         bool
	Err        error
	ResponseTime time.Duration
}

func CheckSite(client *http.Client, url string) Result {

	start := time.Now()

	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return Result{
			URL: url,
			OK:  false,
			Err: err,
			ResponseTime: duration,
		}
	}

	defer resp.Body.Close()

	ok := resp.StatusCode == http.StatusOK

	return Result{
		URL:          url,
		StatusCode:   resp.StatusCode,
		OK:           ok,
		ResponseTime: duration,
	}
}

