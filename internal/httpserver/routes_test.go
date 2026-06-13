package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"serverigh/internal/config"
)

func TestRoutesWireCoreHandlers(t *testing.T) {
	app := New(testConfig(t))

	tests := []struct {
		name   string
		path   string
		status int
	}{
		{
			name:   "root",
			path:   "/",
			status: http.StatusFound,
		},
		{
			name:   "browse",
			path:   "/browse?path=/docs",
			status: http.StatusOK,
		},
		{
			name:   "health",
			path:   "/healthz",
			status: http.StatusOK,
		},
		{
			name:   "unknown",
			path:   "/missing",
			status: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("test request: %v", err)
			}

			if resp.StatusCode != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

func testConfig(t *testing.T) config.Config {
	t.Helper()

	return config.Config{
		Root:            t.TempDir(),
		Host:            "127.0.0.1",
		Port:            4173,
		MaxPreviewBytes: 1024,
	}
}
