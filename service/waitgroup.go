package service

import (
	"apiracer/request"
	"apiracer/utils"
	"fmt"
	"sync"
	"time"
)

func RunWaitGroup(client request.APIClient) (time.Duration, int, utils.ProfilingData) {

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0

	memBefore := utils.CaptureMemStats()
	start := time.Now()

	for i, url := range apiURLs {
		wg.Add(1)
		go func(url string, i int) {
			defer wg.Done()

			resp, err := client.Fetch(url)

			mu.Lock()
			if err != nil {
				fmt.Printf("  [API %d] Failed: %v\n", i+1, err)
			} else {
				successCount++
				fmt.Printf("  [API %d] Success | Location: %v\n", i+1, resp.GeoInfo["Name"])
			}
			mu.Unlock()

		}(url, i)
	}

	goroutines := utils.CaptureGoroutines()
	wg.Wait()
	memAfter := utils.CaptureMemStats()
	duration := time.Since(start)

	profilingData := utils.ProfilingData{
		MemBefore:  memBefore,
		MemAfter:   memAfter,
		Goroutines: goroutines,
		TimeTaken:  duration,
	}
	return duration, successCount, profilingData
}