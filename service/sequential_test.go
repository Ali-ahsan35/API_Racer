package service

import (
    "apiracer/request"
    "errors"
	"sync"
    "testing"

    "github.com/stretchr/testify/assert"
)

// mockClient implements request.APIClient for testing
type mockClient struct {
    response *request.APIResponse
    err      error
}

func (m *mockClient) Fetch(url string) (*request.APIResponse, error) {
    return m.response, m.err
}

func TestRunSequential_AllSuccess(t *testing.T) {
    client := &mockClient{
        response: &request.APIResponse{
            GeoInfo: map[string]interface{}{"Name": "Texas"},
            Result:  map[string]interface{}{"status": "ok"},
            Success: true,
        },
    }

    duration, count, _ := RunSequential(client)

    assert.Equal(t, len(apiURLs), count)   // all 12 succeed
    assert.Positive(t, duration)
}

func TestRunSequential_AllFail(t *testing.T) {
    client := &mockClient{
        err: errors.New("connection refused"),
    }

    _, count, _ := RunSequential(client)

    assert.Equal(t, 0, count)
}

func TestRunSequential_PartialFailure(t *testing.T) {
	var partialClient request.APIClient = &partialMockClient{failAfter: 4}

	_, count, _ := RunSequential(partialClient)

	assert.Equal(t, 4, count)
}
// partialMockClient succeeds for the first N calls, then fails
// in sequential_test.go — replace partialMockClient with this:
type partialMockClient struct {
    mu        sync.Mutex
    calls     int
    failAfter int
}

func (p *partialMockClient) Fetch(url string) (*request.APIResponse, error) {
    p.mu.Lock()
    p.calls++
    calls := p.calls
    p.mu.Unlock()

    if calls > p.failAfter {
        return nil, errors.New("upstream error")
    }
    return &request.APIResponse{
        GeoInfo: map[string]interface{}{"Name": "Hawaii"},
        Result:  map[string]interface{}{},
        Success: true,
    }, nil
}