package checker

import (
	"net/http"
	"time"
)

// Result хранит результат проверки сайта
type Result struct {
	URL        string
	StatusCode int
	OK         bool
	Err        error
}

// CheckSite выполняет HTTP GET запрос и проверяет доступность сайта
func CheckSite(url string) Result {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

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
		Err:        nil,
	}
}
