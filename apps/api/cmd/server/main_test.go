package main

import (
	"bytes"
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"unggulmart/marketing-workspace/apps/api/internal/app"
)

func TestInvalidConfigStopsBeforeListening(t *testing.T) {
	var output bytes.Buffer
	err := run(context.Background(), func(key string) (string, bool) {
		if key == "APP_ENV" {
			return "private-invalid-value", true
		}
		return "", false
	}, &output)
	if err == nil || !strings.Contains(output.String(), "CONFIG_INVALID") || strings.Contains(output.String(), "private") {
		t.Fatalf("startup did not fail safely: %v %s", err, output.String())
	}
}

func TestLocalServerUsesSharedRouterAndBoundedIO(t *testing.T) {
	values := map[string]string{"APP_ENV": "development", "APP_ORIGIN": "http://127.0.0.1:5173"}
	lookup := func(key string) (string, bool) { value, ok := values[key]; return value, ok }
	cfg, err := app.LoadConfig(lookup)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := app.NewRouter(cfg, app.NewLogger(io.Discard, cfg.LogLevel))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server := newHTTPServer(cfg, handler, ctx)
	if server.Addr != cfg.HTTPAddr || server.ReadHeaderTimeout != 5*time.Second ||
		server.WriteTimeout <= cfg.RequestTimeout || server.ReadTimeout <= cfg.RequestTimeout ||
		server.IdleTimeout <= 0 || server.MaxHeaderBytes != 16<<10 || server.BaseContext(nil) != ctx {
		t.Fatal("local server lost HTTP/context limits")
	}
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/tasks/synthetic", nil))
	if w.Code != 404 || !strings.Contains(w.Body.String(), `"code":"NOT_FOUND"`) {
		t.Fatal("router not reused")
	}
	cancel()
	if err := run(ctx, lookup, io.Discard); err != nil {
		t.Fatal("cancelled startup failed unexpectedly")
	}
}
