package service

import (
	"apiracer/request"
	"fmt"
	"sync"
	"time"
)

func RunWaitGroup() (time.Duration, int) {

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0

	start := time.Now()

	for i, url := range apiURLs {
		wg.Add(1)
		go func(url string, i int) {
			defer wg.Done()

			_, err := request.FetchAPI(url)

			mu.Lock()
			if err != nil {
				fmt.Printf("  [API %d] Failed: %v\n", i+1, err)
			} else {
				successCount++
				fmt.Printf("  [API %d] Success\n", i+1)
			}
			mu.Unlock()

		}(url, i)
	}

	wg.Wait()

	duration := time.Since(start)
	return duration, successCount
}