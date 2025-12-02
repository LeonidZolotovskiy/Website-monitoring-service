package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	sites := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.python.org",
		"https://www.nonexistentsite12345.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://www.wikipedia.org",
		"https://www.fakewebsite987654.com",
		"https://www.yahoo.com",
		"https://www.microsoft.com",
	}

	for {
		for _, site := range sites {
			checkSite(site)
		}
		
		time.Sleep(1 * time.Minute)
	}
}

func checkSite(url string) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Site %s NOT ok\n", url)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Site %s ok\n", url)
	} else {
		fmt.Printf("Site %s NOT ok\n", url)
	}
}