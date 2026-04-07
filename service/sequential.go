package service

import (
	"apiracer/request"
	"apiracer/utils"
	"fmt"
	"time"
)

var apiURLs = []string{
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/usa:hawaii?amenities=19&device=desktop&items=1&limit=8&locations=US&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/north-america?device=desktop&items=1&limit=8&locations=US&order=5&pt=9&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/north-america?device=desktop&items=1&limit=8&locations=US&pt=3&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/north-america?device=desktop&items=1&limit=8&locations=US&pt=7&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/north-america?amenities=19-20&device=desktop&items=1&limit=8&locations=US&pt=6&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/usa:texas?device=desktop&items=1&limit=8&locations=US&order=1&pt=11&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/usa:texas?device=desktop&items=1&limit=8&locations=US&order=3&pt=9&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/usa:texas?device=desktop&items=1&limit=8&locations=US&order=1&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/usa:texas?device=desktop&items=1&limit=8&locations=US&nearby=1&order=1&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/usa:texas?amenities=20&device=desktop&items=1&limit=8&locations=US&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/usa:texas?device=desktop&items=1&limit=8&locations=US&order=5&showFallbackData=1",
	"https://ownerdirect.beta.123presto.com/api/v1/category/details/usa:texas?amenities=11&device=desktop&items=1&limit=8&locations=US&showFallbackData=1",
}

func RunSequential(client request.APIClient) (time.Duration, int, utils.ProfilingData) {

	start := time.Now()
	successCount := 0
	memBefore := utils.CaptureMemStats()
	goroutines := utils.CaptureGoroutines()

	for i, url := range apiURLs {
		resp, err := client.Fetch(url)
		if err != nil {
			fmt.Printf("  [API %d] Failed: %v\n", i+1, err)
		} else {
			successCount++
			fmt.Printf("  [API %d] Success | Location: %v\n", i+1, resp.GeoInfo["Name"])
		}
	}

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