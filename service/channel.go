package service

import (
	"apiracer/request"
	"fmt"
	"sync"
	"time"
)

func RunChannel() (time.Duration, int) {

	ch := make(chan bool, len(apiURLs))
	var wg sync.WaitGroup
	start := time.Now()

	for i, url := range apiURLs {
		wg.Add(1)
		go func(url string, i int) {
			defer wg.Done()
			resp, err := request.FetchAPI(url)
			if err != nil {
				fmt.Printf("  [API %d] Failed: %v\n", i+1, err)
				ch <- false
			} else {
				fmt.Printf("  [API %d] Success | Location: %v\n", i+1, resp.GeoInfo["Name"])
				ch <- true
			}
		}(url, i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	successCount := 0
	for result := range ch {
		if result {
			successCount++
		}
	}

	duration := time.Since(start)
	return duration, successCount
}
