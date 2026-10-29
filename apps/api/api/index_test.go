package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"unggulmart/marketing-workspace/apps/api/internal/platform/httpx"
)

func TestExportedHandlerUsesSharedHTTPBoundary(t *testing.T) {
	// Synthetic process-local environment; never reads a database or sends traffic.
	for key, value := range map[string]string{
		"APP_ENV": "development", "VERCEL_ENV": "development",
		"APP_ORIGIN": "http://127.0.0.1:5173", "HTTP_ADDR": "127.0.0.1:8080",
		"REQUEST_TIMEOUT": "1s", "LOG_LEVEL": "error", "PUSH_ENABLED": "false", "CRON_ENABLED": "false",
	} {
		t.Setenv(key, value)
	}
	for _, path := range []string{"/", "/api/v1/tasks/synthetic?filter=synthetic"} {
		w := httptest.NewRecorder()
		Handler(w, httptest.NewRequest("GET", path, nil))
		var body httpx.ErrorEnvelope
		if w.Code != 404 || json.Unmarshal(w.Body.Bytes(), &body) != nil || body.Error.Code != httpx.NotFound ||
			body.Error.RequestID != w.Header().Get("X-Request-ID") || body.Error.RequestID == "" {
			t.Fatalf("adapter lost shared error contract: %d %s", w.Code, w.Body.String())
		}
	}
	w := httptest.NewRecorder()
	Handler(w, httptest.NewRequest("POST", "/api/v1/tasks/synthetic", nil))
	if w.Code != 403 {
		t.Fatalf("adapter bypassed Origin guard: %d", w.Code)
	}
}
