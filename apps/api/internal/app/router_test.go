package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"unggulmart/marketing-workspace/apps/api/internal/platform/httpx"
)

func TestRouterPreservesNestedRequestAndResponse(t *testing.T) {
	cfg, _ := LoadConfig(lookupMap(devEnv()))
	var logs bytes.Buffer
	router, err := NewRouter(cfg, NewLogger(&logs, cfg.LogLevel))
	if err != nil {
		t.Fatal(err)
	}
	// Fixture endpoint only; never mounted in the runtime application.
	router.Post("/api/v1/fixture/{id}", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Title string `json:"title"`
		}
		if code := httpx.DecodeJSON(w, r, &input); code != "" {
			httpx.WriteError(w, r, code)
			return
		}
		if chi.URLParam(r, "id") != "fixture-id" || input.Title != "synthetic title" ||
			r.URL.Query().Get("filter") != "synthetic" || r.Method != "POST" ||
			r.Header.Get("Cookie") != "synthetic_cookie=private-cookie" {
			t.Error("request path/method/query/body/header was changed")
		}
		w.Header().Add("Set-Cookie", "first=synthetic; HttpOnly")
		w.Header().Add("Set-Cookie", "second=synthetic; HttpOnly")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{"data":{"ok":true}}`)
	})
	req := httptest.NewRequest("POST", "/api/v1/fixture/fixture-id?filter=synthetic&token=private-query", strings.NewReader(`{"title":"synthetic title"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", cfg.AppOrigin)
	req.Header.Set("Cookie", "synthetic_cookie=private-cookie")
	req.Header.Set("Authorization", "Bearer private-authorization")
	req.Header.Set("X-CSRF-Token", "private-csrf")
	req.Header.Set("X-Request-ID", "private-forged-request-id")
	req.Header.Set("X-Forwarded-For", "private-forwarded-ip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 201 || len(w.Header().Values("Set-Cookie")) != 2 || w.Body.String() != `{"data":{"ok":true}}` {
		t.Fatalf("response changed: %d %v %s", w.Code, w.Header(), w.Body.String())
	}
	if w.Header().Get("X-Request-ID") == "private-forged-request-id" || !strings.HasPrefix(w.Header().Get("X-Request-ID"), "req_") {
		t.Fatal("untrusted request ID accepted")
	}
	if strings.Contains(logs.String(), "private") || strings.Contains(logs.String(), "synthetic title") {
		t.Fatalf("request data leaked to log: %s", logs.String())
	}
	var entry map[string]any
	if json.Unmarshal(logs.Bytes(), &entry) != nil || entry["status"] != float64(201) ||
		entry["request_id"] != w.Header().Get("X-Request-ID") || entry["method"] != "POST" {
		t.Fatalf("missing safe log metadata: %s", logs.String())
	}
}

func TestErrorsAndOriginPolicy(t *testing.T) {
	cfg, _ := LoadConfig(lookupMap(devEnv()))
	router, _ := NewRouter(cfg, NewLogger(io.Discard, cfg.LogLevel))
	router.Get("/api/v1/fixture", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })
	cases := []struct {
		method, path, origin string
		status               int
		code                 httpx.ErrorCode
	}{
		{"GET", "/api/v1/missing/private", "", 404, httpx.NotFound},
		{"POST", "/api/v1/fixture", cfg.AppOrigin, 405, httpx.MethodNotAllowed},
		{"POST", "/api/v1/fixture", "", 403, httpx.InvalidOrigin},
		{"PATCH", "/api/v1/missing", "https://other.example.invalid", 403, httpx.InvalidOrigin},
		{"PUT", "/api/v1/missing", "null", 403, httpx.InvalidOrigin},
		{"DELETE", "/api/v1/missing", cfg.AppOrigin + "/", 403, httpx.InvalidOrigin},
	}
	for _, test := range cases {
		t.Run(test.method+test.path+test.origin, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.path, nil)
			if test.origin != "" {
				req.Header.Set("Origin", test.origin)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assertError(t, w, test.status, test.code)
			if test.status == 405 && w.Header().Get("Allow") != "GET" {
				t.Fatal("405 must include the permitted method")
			}
		})
	}
	req := httptest.NewRequest("POST", "/api/v1/fixture", nil)
	req.Header.Add("Origin", cfg.AppOrigin)
	req.Header.Add("Origin", cfg.AppOrigin)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertError(t, w, 403, httpx.InvalidOrigin)
}

func TestNestedMethodNotAllowedPreservesAllowHeader(t *testing.T) {
	cfg, _ := LoadConfig(lookupMap(devEnv()))
	router, _ := NewRouter(cfg, NewLogger(io.Discard, cfg.LogLevel))
	router.Route("/api/v1/nested", func(child chi.Router) {
		child.Get("/{id}", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })
	})
	req := httptest.NewRequest("POST", "/api/v1/nested/synthetic", nil)
	req.Header.Set("Origin", cfg.AppOrigin)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assertError(t, w, 405, httpx.MethodNotAllowed)
	if w.Header().Get("Allow") != "GET" {
		t.Fatalf("nested Allow header lost: %v", w.Header())
	}
}

func TestRecoveryDiscardsPartialSuccessAndPanicDetails(t *testing.T) {
	cfg, _ := LoadConfig(lookupMap(devEnv()))
	var logs bytes.Buffer
	router, _ := NewRouter(cfg, NewLogger(&logs, cfg.LogLevel))
	router.Get("/api/v1/fixture", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "private-response-cookie")
		w.WriteHeader(200)
		_, _ = io.WriteString(w, "private-partial-body")
		panic("private-database-stack-token")
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/fixture?secret=private", nil))
	assertError(t, w, 500, httpx.InternalError)
	if strings.Contains(logs.String(), "private") || strings.Contains(w.Body.String(), "private") || w.Header().Get("Set-Cookie") != "" {
		t.Fatalf("panic/partial response leaked: %s %s", logs.String(), w.Body.String())
	}
}

func TestInvalidConfigReturnsUnavailableWithoutEchoingValues(t *testing.T) {
	var logs bytes.Buffer
	handler := FromEnvironment(lookupMap(map[string]string{"APP_ENV": "private-invalid-config"}), &logs)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/tasks/private?token=private", nil))
	assertError(t, w, 503, httpx.ServiceUnavailable)
	if strings.Contains(logs.String(), "private") || strings.Contains(w.Body.String(), "private") {
		t.Fatal("configuration or request leaked")
	}
	if !strings.Contains(logs.String(), "CONFIG_INVALID") {
		t.Fatal("missing safe startup error")
	}
}

func TestSecurityHeadersFollowConfiguredEnvironment(t *testing.T) {
	for _, env := range []string{"development", "preview", "production"} {
		values := devEnv()
		values["APP_ENV"] = env
		if env != "development" {
			values["APP_ORIGIN"] = "https://web.example.invalid"
		}
		handler := FromEnvironment(lookupMap(values), io.Discard)
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Forwarded-Proto", "https")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if (w.Header().Get("Strict-Transport-Security") != "") != (env == "production") {
			t.Fatalf("HSTS must follow config, not forwarded headers: %s", env)
		}
		if w.Header().Get("X-Content-Type-Options") != "nosniff" || w.Header().Get("Content-Security-Policy") == "" {
			t.Fatal("missing security headers")
		}
	}
}

func TestRouterConstructorRejectsInvalidBoundary(t *testing.T) {
	cfg, _ := LoadConfig(lookupMap(devEnv()))
	if _, err := NewRouter(cfg, nil); err == nil {
		t.Fatal("nil logger accepted")
	}
	cfg.Environment = ""
	if _, err := NewRouter(cfg, NewLogger(io.Discard, 0)); err == nil {
		t.Fatal("invalid configuration accepted")
	}
}

func assertError(t *testing.T, w *httptest.ResponseRecorder, status int, code httpx.ErrorCode) {
	t.Helper()
	var body httpx.ErrorEnvelope
	if w.Code != status || json.Unmarshal(w.Body.Bytes(), &body) != nil || body.Error.Code != code {
		t.Fatalf("unexpected error: %d %s", w.Code, w.Body.String())
	}
	if body.Error.RequestID == "" || body.Error.RequestID != w.Header().Get("X-Request-ID") || body.Error.Fields == nil ||
		w.Header().Get("Content-Type") != "application/json; charset=utf-8" || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("error contract incomplete: %v %s", w.Header(), w.Body.String())
	}
}
