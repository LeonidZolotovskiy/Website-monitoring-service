package main

import (
	"fmt"
	"site-monitor/internal/checker"
)

func main() {
	fmt.Println("Site Monitor started")

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

	for _, site := range sites {
		checker.CheckSite(site)
	}
}