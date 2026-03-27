package service

import (
	"apiracer/request"
	"fmt"
	"time"
)

func RunChannel() (time.Duration, int) {

	ch := make(chan bool, len(apiURLs))
	start := time.Now()

	for i, url := range apiURLs {
		go func(url string, i int) {
			_, err := request.FetchAPI(url)
			if err != nil {
				fmt.Printf("  [API %d] Failed\n", i+1)
				ch <- false
			} else {
				fmt.Printf("  [API %d] Success\n", i+1)
				ch <- true
			}
		}(url, i)
	}

	successCount := 0
	for i := 0; i < len(apiURLs); i++ {
		result := <-ch
		if result {
			successCount++
		}
	}

	duration := time.Since(start)

	return duration, successCount
}
