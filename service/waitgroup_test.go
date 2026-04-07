package service

import (
	"apiracer/request"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunWaitGroup_AllSuccess(t *testing.T) {
	client := &mockClient{
		response: &request.APIResponse{
			GeoInfo: map[string]interface{}{"Name": "Texas"},
			Result:  map[string]interface{}{"status": "ok"},
			Success: true,
		},
	}

	duration, count, _ := RunWaitGroup(client)

	assert.Equal(t, len(apiURLs), count)
	assert.Positive(t, duration)
}

func TestRunWaitGroup_AllFail(t *testing.T) {
	client := &mockClient{
		err: errors.New("connection refused"),
	}

	_, count, _ := RunWaitGroup(client)

	assert.Equal(t, 0, count)
}

func TestRunWaitGroup_PartialFailure(t *testing.T) {
	client := &partialMockClient{failAfter: 6}

	_, count, _ := RunWaitGroup(client)

	assert.Equal(t, 6, count)
}

func TestRunWaitGroup_ConcurrentSafe(t *testing.T) {
	// Run multiple times to surface any race conditions
	// Always run this with: go test -race ./service/...
	client := &mockClient{
		response: &request.APIResponse{
			GeoInfo: map[string]interface{}{"Name": "Hawaii"},
			Result:  map[string]interface{}{},
			Success: true,
		},
	}

	for i := 0; i < 5; i++ {
		_, count, _ := RunWaitGroup(client)
		assert.Equal(t, len(apiURLs), count)
	}
}

func TestRunWaitGroup_ProfilingDataReturned(t *testing.T) {
	client := &mockClient{
		response: &request.APIResponse{
			GeoInfo: map[string]interface{}{"Name": "Texas"},
			Result:  map[string]interface{}{},
			Success: true,
		},
	}

	duration, _, profiling := RunWaitGroup(client)

	assert.Positive(t, duration)
	assert.Positive(t, profiling.Goroutines)
	assert.Positive(t, profiling.TimeTaken)
}