package service

import (
	"apiracer/request"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunChannel_AllSuccess(t *testing.T) {
	client := &mockClient{
		response: &request.APIResponse{
			GeoInfo: map[string]interface{}{"Name": "Texas"},
			Result:  map[string]interface{}{"status": "ok"},
			Success: true,
		},
	}

	duration, count, _ := RunChannel(client)

	assert.Equal(t, len(apiURLs), count)
	assert.Positive(t, duration)
}

func TestRunChannel_AllFail(t *testing.T) {
	client := &mockClient{
		err: errors.New("connection refused"),
	}

	_, count, _ := RunChannel(client)

	assert.Equal(t, 0, count)
}

func TestRunChannel_PartialFailure(t *testing.T) {
	client := &partialMockClient{failAfter: 7}

	_, count, _ := RunChannel(client)

	assert.Equal(t, 7, count)
}

func TestRunChannel_ChannelDrainsCompletely(t *testing.T) {
	// Verifies every goroutine sends exactly one result into the channel
	// and the channel closes cleanly — no deadlock, no leak
	client := &mockClient{
		response: &request.APIResponse{
			GeoInfo: map[string]interface{}{"Name": "Hawaii"},
			Result:  map[string]interface{}{},
			Success: true,
		},
	}

	_, count, _ := RunChannel(client)

	// All 12 results must be collected — if channel leaked, count would be wrong
	assert.Equal(t, len(apiURLs), count)
}

func TestRunChannel_ConcurrentSafe(t *testing.T) {
	// Always run with: go test -race ./service/...
	client := &mockClient{
		response: &request.APIResponse{
			GeoInfo: map[string]interface{}{"Name": "Texas"},
			Result:  map[string]interface{}{},
			Success: true,
		},
	}

	for i := 0; i < 5; i++ {
		_, count, _ := RunChannel(client)
		assert.Equal(t, len(apiURLs), count)
	}
}

func TestRunChannel_ProfilingDataReturned(t *testing.T) {
	client := &mockClient{
		response: &request.APIResponse{
			GeoInfo: map[string]interface{}{"Name": "Texas"},
			Result:  map[string]interface{}{},
			Success: true,
		},
	}

	duration, _, profiling := RunChannel(client)

	assert.Positive(t, duration)
	assert.Positive(t, profiling.Goroutines)
	assert.Positive(t, profiling.TimeTaken)
}