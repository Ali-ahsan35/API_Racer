package request

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setAPIKey sets or clears the beego apikey config for a test.
func setAPIKey(t *testing.T, key string) {
	t.Helper()
	beego.AppConfig.Set("apikey", key)
	t.Cleanup(func() {
		beego.AppConfig.Set("apikey", "")
	})
}

// validPayload is a minimal response that satisfies all three validation levels.
const validPayload = `{"GeoInfo": {"country": "USA", "city": "Hawaii"}, "Result": {"status": "ok"}}`

// ─── Happy path ────────────────────────────────────────────────────────────────

func TestFetchAPI_Success(t *testing.T) {
	setAPIKey(t, "test_key")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify all three required headers are forwarded
		assert.Equal(t, "test_key", r.Header.Get("x-api-key"))
		assert.Equal(t, "XMLHttpRequest", r.Header.Get("X-Requested-With"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validPayload))
	}))
	defer server.Close()

	resp, err := FetchAPI(server.URL)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "USA", resp.GeoInfo["country"])
	assert.Equal(t, "ok", resp.Result["status"])
	assert.Equal(t, server.URL, resp.URL)  // URL must be stored on the response
	assert.Empty(t, resp.Error)
}

// ─── API key validation ────────────────────────────────────────────────────────

func TestFetchAPI_MissingAPIKey(t *testing.T) {
	// Ensure the key is definitely absent
	beego.AppConfig.Set("apikey", "")

	resp, err := FetchAPI("http://does-not-matter.example.com")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "API key not found")
}

// ─── HTTP status codes ────────────────────────────────────────────────────────

func TestFetchAPI_NonOKStatusCodes(t *testing.T) {
	setAPIKey(t, "test_key")

	cases := []struct {
		name   string
		status int
	}{
		{"internal server error", http.StatusInternalServerError},
		{"not found", http.StatusNotFound},
		{"unauthorized", http.StatusUnauthorized},
		{"bad gateway", http.StatusBadGateway},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
			}))
			defer server.Close()

			resp, err := FetchAPI(server.URL)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), "unexpected status code")
		})
	}
}

// ─── JSON parsing ─────────────────────────────────────────────────────────────

func TestFetchAPI_JSONErrors(t *testing.T) {
	setAPIKey(t, "test_key")

	cases := []struct {
		name        string
		body        string
		errContains string
	}{
		{
			name:        "completely invalid JSON",
			body:        `invalid json`,
			errContains: "invalid JSON",
		},
		{
			name:        "empty body",
			body:        ``,
			errContains: "invalid JSON",
		},
		{
			name:        "JSON array instead of object",
			body:        `[1, 2, 3]`,
			errContains: "invalid JSON",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tc.body))
			}))
			defer server.Close()

			resp, err := FetchAPI(server.URL)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tc.errContains)
		})
	}
}

// ─── Missing required fields ───────────────────────────────────────────────────

func TestFetchAPI_MissingFields(t *testing.T) {
	setAPIKey(t, "test_key")

	cases := []struct {
		name        string
		body        string
		errContains string
	}{
		{
			name:        "missing GeoInfo",
			body:        `{"Result": {"status": "ok"}}`,
			errContains: "missing GeoInfo",
		},
		{
			name:        "missing Result",
			body:        `{"GeoInfo": {"country": "USA"}}`,
			errContains: "missing Result",
		},
		{
			name:        "empty JSON object — both fields missing",
			body:        `{}`,
			errContains: "missing GeoInfo",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tc.body))
			}))
			defer server.Close()

			resp, err := FetchAPI(server.URL)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tc.errContains)
		})
	}
}

// ─── Type safety — panic guard ────────────────────────────────────────────────
// FetchAPI uses a bare type assertion: geoInfo.(map[string]interface{})
// If the API returns a non-object for GeoInfo or Result, this panics.
// These tests document that known risk and will fail if the code is not guarded.

func TestFetchAPI_GeoInfoWrongType_ShouldNotPanic(t *testing.T) {
	setAPIKey(t, "test_key")

	// GeoInfo is a string, not an object — bare assertion will panic
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"GeoInfo": "not-a-map", "Result": {"status": "ok"}}`))
	}))
	defer server.Close()

	// If the assertion is unguarded, this will panic and the test runner marks it as failed.
	assert.NotPanics(t, func() {
		resp, err := FetchAPI(server.URL)
		// After fixing the production code, expect an error, not a panic
		_ = resp
		_ = err
	})
}

func TestFetchAPI_ResultWrongType_ShouldNotPanic(t *testing.T) {
	setAPIKey(t, "test_key")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"GeoInfo": {"country": "USA"}, "Result": 42}`))
	}))
	defer server.Close()

	assert.NotPanics(t, func() {
		resp, err := FetchAPI(server.URL)
		_ = resp
		_ = err
	})
}

// ─── Network errors ───────────────────────────────────────────────────────────

func TestFetchAPI_InvalidURL(t *testing.T) {
	setAPIKey(t, "test_key")

	// An invalid URL fails at http.NewRequest, before any network call
	resp, err := FetchAPI("://bad-url")

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestFetchAPI_ServerUnreachable(t *testing.T) {
	setAPIKey(t, "test_key")

	// Port 1 is almost universally refused immediately
	resp, err := FetchAPI("http://127.0.0.1:1/api")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "request failed")
}

// ─── Response struct completeness ─────────────────────────────────────────────

func TestFetchAPI_ResponseFieldsPopulated(t *testing.T) {
	setAPIKey(t, "test_key")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"GeoInfo": {"country": "Bangladesh", "city": "Dhaka"},
			"Result":  {"ip": "1.2.3.4", "isp": "TestISP"}
		}`))
	}))
	defer server.Close()

	resp, err := FetchAPI(server.URL)

	require.NoError(t, err)
	require.NotNil(t, resp)

	// Every field on APIResponse must be set correctly
	assert.True(t, resp.Success)
	assert.Empty(t, resp.Error)
	assert.Equal(t, server.URL, resp.URL)
	assert.Equal(t, "Bangladesh", resp.GeoInfo["country"])
	assert.Equal(t, "Dhaka", resp.GeoInfo["city"])
	assert.Equal(t, "1.2.3.4", resp.Result["ip"])
	assert.Equal(t, "TestISP", resp.Result["isp"])
}
// Make timeout configurable
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
}

// In tests you can override it
func TestFetchAPI_Timeout(t *testing.T) {
    // Override client with 1 second timeout for testing
    original := httpClient
    httpClient = &http.Client{Timeout: 1 * time.Second}
    defer func() { httpClient = original }()

    setAPIKey(t, "test_key")

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(3 * time.Second) // only needs 3s now!
        w.WriteHeader(http.StatusOK)
    }))
    defer server.Close()

    resp, err := FetchAPI(server.URL)

    assert.Error(t, err)
    assert.Nil(t, resp)
    assert.Contains(t, err.Error(), "request failed")
}