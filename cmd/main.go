package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"site-monitor/internal/checker"
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

	
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigCh
		fmt.Println("\nReceived signal:", sig)
		cancel()
	}()

	fmt.Println("Starting site checks... Press Ctrl+C to stop.")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down gracefully...")
			return
		default:
			for _, site := range sites {
				checker.CheckSite(site)
			}
			time.Sleep(1 * time.Minute)
		}
	}
}