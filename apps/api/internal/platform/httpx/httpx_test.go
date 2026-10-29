package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func pipeline(handler http.Handler, timeout time.Duration) http.Handler {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return RequestIDs(Headers(false)(Recoverer(logger)(Deadline(timeout)(LimitBody(handler)))))
}

func checkError(t *testing.T, w *httptest.ResponseRecorder, code ErrorCode, status int) {
	t.Helper()
	var body ErrorEnvelope
	if w.Code != status || json.Unmarshal(w.Body.Bytes(), &body) != nil || body.Error.Code != code ||
		body.Error.RequestID == "" || body.Error.RequestID != w.Header().Get("X-Request-ID") || body.Error.Fields == nil {
		t.Fatalf("incorrect error: %d %s", w.Code, w.Body.String())
	}
	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("missing JSON/no-store headers")
	}
}

func TestDeadlineCancelsWorkAndRejectsLateWrites(t *testing.T) {
	release := make(chan struct{})
	cancelled := make(chan struct{})
	writeResult := make(chan error, 1)
	handler := pipeline(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		close(cancelled)
		<-release
		_, err := io.WriteString(w, "private-late-write")
		writeResult <- err
	}), 20*time.Millisecond)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	checkError(t, w, RequestTimeout, 503)
	select {
	case <-cancelled:
	case <-time.After(time.Second):
		t.Fatal("handler did not receive cancellation")
	}
	close(release)
	select {
	case err := <-writeResult:
		if !errors.Is(err, http.ErrHandlerTimeout) {
			t.Fatalf("late write accepted: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("handler failed to exit")
	}
	if strings.Contains(w.Body.String(), "private") {
		t.Fatal("late write changed response")
	}
}

func TestParentCancellationNeverPublishesSuccess(t *testing.T) {
	for range 32 {
		ctx, cancel := context.WithCancel(context.Background())
		handler := pipeline(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cancel()
			w.Header().Set("Set-Cookie", "private-late-cookie")
			w.WriteHeader(200)
			_, _ = io.WriteString(w, "private-late-success")
		}), time.Second)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil).WithContext(ctx))
		cancel()
		checkError(t, w, RequestCancelled, 503)
		if strings.Contains(w.Body.String(), "private") || w.Header().Get("Set-Cookie") != "" {
			t.Fatal("cancelled request published a partial response")
		}
	}
}

func TestAlreadyCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	handler := pipeline(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !errors.Is(r.Context().Err(), context.Canceled) {
			t.Error("missing cancelled parent context")
		}
	}), time.Second)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil).WithContext(ctx))
	checkError(t, w, RequestCancelled, 503)
}

func TestJSONBoundary(t *testing.T) {
	handler := pipeline(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var value struct {
			Title string `json:"title"`
		}
		if code := DecodeJSON(w, r, &value); code != "" {
			WriteError(w, r, code)
			return
		}
		w.WriteHeader(204)
	}), time.Second)
	cases := []struct {
		name, body, contentType string
		chunked                 bool
		status                  int
		code                    ErrorCode
	}{
		{"valid", `{"title":"synthetic"}`, "application/json; charset=utf-8", false, 204, ""},
		{"unicode", `{"title":"Judul üji 🎉"}`, "application/json", false, 204, ""},
		{"invalid utf8", "{\"title\":\"" + string([]byte{0xff}) + "\"}", "application/json", false, 400, InvalidJSON},
		{"unknown", `{"title":"ok","completed_by":"private"}`, "application/json", false, 400, InvalidJSON},
		{"broken", `{"title":`, "application/json", false, 400, InvalidJSON},
		{"empty", ``, "application/json", false, 400, InvalidJSON},
		{"multiple", `{"title":"one"}{"title":"two"}`, "application/json", false, 400, InvalidJSON},
		{"trailing", `{"title":"ok"} private`, "application/json", false, 400, InvalidJSON},
		{"media", `{"title":"ok"}`, "text/plain", false, 415, UnsupportedMedia},
		{"missing media", `{}`, "", false, 415, UnsupportedMedia},
		{"oversized length", `{"title":"` + strings.Repeat("x", int(MaxBodyBytes)) + `"}`, "application/json", false, 413, RequestTooLarge},
		{"oversized chunked", `{"title":"` + strings.Repeat("x", int(MaxBodyBytes)) + `"}`, "application/json", true, 413, RequestTooLarge},
		{"oversized trailing whitespace", `{"title":"ok"}` + strings.Repeat(" ", int(MaxBodyBytes)), "application/json", true, 413, RequestTooLarge},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", strings.NewReader(test.body))
			req.Header.Set("Content-Type", test.contentType)
			if test.chunked {
				req.ContentLength = -1
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			if test.status == 204 {
				if w.Code != 204 {
					t.Fatalf("valid JSON rejected: %d %s", w.Code, w.Body.String())
				}
			} else {
				checkError(t, w, test.code, test.status)
				if strings.Contains(w.Body.String(), "private") {
					t.Fatal("decoder input leaked")
				}
			}
		})
	}
}

func TestRequestIDsConcurrentAndUntrustedMethodRedaction(t *testing.T) {
	var wg sync.WaitGroup
	ids := make(chan string, 64)
	handler := RequestIDs(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if RequestID(r) != w.Header().Get("X-Request-ID") {
			t.Error("request ID mismatch")
		}
		ids <- RequestID(r)
	}))
	for range 64 {
		wg.Go(func() {
			req := httptest.NewRequest("GET", "/private", nil)
			req.Header.Set("X-Request-ID", "private-forged-id")
			handler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
	wg.Wait()
	close(ids)
	seen := map[string]bool{}
	for id := range ids {
		if !strings.HasPrefix(id, "req_") || seen[id] {
			t.Fatal("duplicate/invalid request ID")
		}
		seen[id] = true
	}
	var logs bytes.Buffer
	logged := RequestIDs(AccessLog(slog.New(slog.NewJSONHandler(&logs, nil)))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))
	logged.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("private-method", "/private?token=private", nil))
	if strings.Contains(logs.String(), "private") || !strings.Contains(logs.String(), `"method":"OTHER"`) {
		t.Fatalf("untrusted data entered log: %s", logs.String())
	}
}

func TestUnknownErrorCodeDoesNotLeak(t *testing.T) {
	handler := RequestIDs(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, r, ErrorCode("private-error-details"))
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	checkError(t, w, InternalError, 500)
	if strings.Contains(w.Body.String(), "private") {
		t.Fatal("unknown error code leaked")
	}
}
