package httpx

import (
	"context"
	"crypto/rand"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type requestIDKey struct{}
type requestStopped struct{ code ErrorCode }

func RequestID(r *http.Request) string {
	value, _ := r.Context().Value(requestIDKey{}).(string)
	return value
}

func RequestIDs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Header dari browser/proxy belum merupakan identitas terpercaya.
		id := "req_" + rand.Text()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id)))
	})
}

func Headers(production bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Referrer-Policy", "no-referrer")
			w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'")
			if production {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Origin(allowed string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				values := r.Header.Values("Origin")
				if len(values) != 1 || values[0] != allowed {
					WriteError(w, r, InvalidOrigin)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AccessLog(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(wrapped, r)
			status := wrapped.Status()
			if status == 0 {
				status = http.StatusOK
			}
			level := slog.LevelInfo
			errorCode := ""
			if status >= 400 {
				errorCode = "HTTP_" + strconv.Itoa(status)
			}
			if status >= 500 {
				level = slog.LevelError
			}
			method := r.Method
			switch method {
			case "GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS":
			default:
				method = "OTHER"
			}
			// Allowlist metadata. Jangan log URL/path/query, IP/forwarded headers,
			// cookies, request/response body, credential, atau arbitrary errors.
			logger.Log(r.Context(), level, "http_request",
				"request_id", RequestID(r), "action", "http_request", "method", method,
				"status", status, "error_code", errorCode, "latency_ms", time.Since(started).Milliseconds())
		})
	}
}

func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recover() != nil {
					logger.ErrorContext(r.Context(), "request_panic", "request_id", RequestID(r), "code", InternalError)
					WriteError(w, r, InternalError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Deadline(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if value := recover(); value != nil {
					if stopped, ok := value.(requestStopped); ok {
						WriteError(w, r, stopped.code)
						return
					}
					panic(value)
				}
			}()
			wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			// net/http buffers normal responses; late writes cannot alter the
			// committed timeout response. Recovery sits OUTSIDE this buffer, so
			// even a handler that wrote partial success then panics returns JSON 500.
			guarded := http.HandlerFunc(func(inner http.ResponseWriter, request *http.Request) {
				next.ServeHTTP(inner, request)
				// Completion and cancellation can both be ready in TimeoutHandler's
				// select. Never publish a success written after cancellation.
				switch request.Context().Err() {
				case context.DeadlineExceeded:
					panic(requestStopped{code: RequestTimeout})
				case context.Canceled:
					panic(requestStopped{code: RequestCancelled})
				}
			})
			http.TimeoutHandler(guarded, timeout, string(errorBody(r, RequestTimeout))).ServeHTTP(wrapped, r)
			if wrapped.Status() == http.StatusServiceUnavailable && wrapped.BytesWritten() == 0 && r.Context().Err() != nil {
				// TimeoutHandler sends an empty 503 for parent cancellation.
				// Keep the same structured error contract without writing headers twice.
				_, _ = wrapped.Write(errorBody(r, RequestCancelled))
			}
		})
	}
}
