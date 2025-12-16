package checker

import (
	"fmt"
	"net/http"
	"time"
)

func CheckSite(url string) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Site %s NOT ok\n", url)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Site %s ok\n", url)
	} else {
		fmt.Printf("Site %s NOT ok (status %d)\n", url, resp.StatusCode)
	}
}