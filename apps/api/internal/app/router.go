package app

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"unggulmart/marketing-workspace/apps/api/internal/platform/httpx"
)

func NewLogger(output io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: level}))
}

type Router struct {
	*chi.Mux
	handler http.Handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	r.handler.ServeHTTP(w, request)
}

// NewRouter belum memasang endpoint domain atau membuka database.
func NewRouter(cfg Config, logger *slog.Logger) (*Router, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if logger == nil {
		return nil, errors.New("logger wajib tersedia")
	}
	router := chi.NewRouter()
	router.NotFound(func(w http.ResponseWriter, r *http.Request) { httpx.WriteError(w, r, httpx.NotFound) })
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.URL.RawPath != "" {
			path = r.URL.RawPath
		}
		allowed := []string{}
		for _, method := range []string{"CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "TRACE"} {
			// Match only inspects routing data, without executing handlers. The
			// root mux also resolves nested routes against the original URL.
			if router.Match(chi.NewRouteContext(), method, path) {
				allowed = append(allowed, method)
			}
		}
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		httpx.WriteError(w, r, httpx.MethodNotAllowed)
	})
	// Chi's empty-mux fast path bypasses Mux.Use. The HTTP boundary must wrap
	// the mux itself so even an empty router/404 receives the same protection.
	handler := chi.Chain(httpx.RequestIDs, httpx.Headers(cfg.Environment == "production"),
		httpx.AccessLog(logger), httpx.Recoverer(logger), httpx.Deadline(cfg.RequestTimeout),
		httpx.Origin(cfg.AppOrigin), httpx.LimitBody).Handler(router)
	return &Router{Mux: router, handler: handler}, nil
}

// FromEnvironment adalah factory cold start adapter serverless. Config invalid
// menghasilkan handler 503, tanpa fallback ke development atau panic ke provider.
func FromEnvironment(lookup Lookup, output io.Writer) http.Handler {
	cfg, err := LoadConfig(lookup)
	if err == nil {
		logger := NewLogger(output, cfg.LogLevel)
		router, routerErr := NewRouter(cfg, logger)
		if routerErr == nil {
			return router
		}
		err = routerErr
	}
	logger := NewLogger(output, slog.LevelInfo)
	logger.Error("startup_failed", "code", "CONFIG_INVALID", "field", configField(err))
	unavailable := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, r, httpx.ServiceUnavailable)
	})
	return httpx.RequestIDs(httpx.Headers(false)(httpx.AccessLog(logger)(unavailable)))
}
