package checker

import (
	"net/http"
)

type Result struct {
	URL        string
	StatusCode int
	OK         bool
	Err        error
}

func CheckSite(client *http.Client, url string) Result {
	resp, err := client.Get(url)
	if err != nil {
		return Result{
			URL: url,
			OK:  false,
			Err: err,
		}
	}
	defer resp.Body.Close()

	ok := resp.StatusCode == http.StatusOK

	return Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		OK:         ok,
	}
}
