package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	
	beego "github.com/beego/beego/v2/server/web"
)

type APIResponse struct {
	GeoInfo map[string]interface{}
	Result  map[string]interface{}
	Success bool
	Error   string
	URL     string
}

func FetchAPI(url string) (*APIResponse, error) {

	apiKey, err := beego.AppConfig.String("apikey")
	if err != nil || apiKey == "" {
		return nil, fmt.Errorf("API key not found in config")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var rawResult map[string]interface{}
	err = json.Unmarshal(body, &rawResult)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON response: %v", err)
	}

	geoInfo, hasGeoInfo := rawResult["GeoInfo"]
	if !hasGeoInfo {
		return nil, fmt.Errorf("missing GeoInfo field in response")
	}

	result, hasResult := rawResult["Result"]
	if !hasResult {
		return nil, fmt.Errorf("missing Result field in response")
	}

	geoMap, ok := geoInfo.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("GeoInfo field is not an object")
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("result field is not an object")
	}

	

	apiResponse := &APIResponse{
		GeoInfo: geoMap,
		Result:  resultMap,
		Success: true,
		Error:   "",
		URL:     url,
	}

	return apiResponse, nil
}