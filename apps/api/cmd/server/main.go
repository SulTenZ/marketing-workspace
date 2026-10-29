// Command server menjalankan router yang sama dengan adapter serverless.
package main

import (
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"unggulmart/marketing-workspace/apps/api/internal/app"
)

func main() {
	if exitCode() != 0 {
		os.Exit(1)
	}
}

func exitCode() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if run(ctx, os.LookupEnv, os.Stdout) != nil {
		return 1
	}
	return 0
}

func newHTTPServer(cfg app.Config, handler http.Handler, ctx context.Context) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       cfg.RequestTimeout + 5*time.Second,
		WriteTimeout:      cfg.RequestTimeout + 5*time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
		BaseContext:       func(net.Listener) context.Context { return ctx },
		// Raw net/http diagnostics can contain user input/stack traces. Application
		// recovery/access logs and lifecycle logs below use safe structured fields.
		ErrorLog: log.New(io.Discard, "", 0),
	}
}

func run(parent context.Context, lookup app.Lookup, output io.Writer) error {
	cfg, err := app.LoadConfig(lookup)
	if err != nil {
		app.NewLogger(output, slog.LevelInfo).Error("startup_failed", "code", "CONFIG_INVALID")
		return err
	}
	logger := app.NewLogger(output, cfg.LogLevel)
	router, err := app.NewRouter(cfg, logger)
	if err != nil {
		logger.Error("startup_failed", "code", "CONFIG_INVALID")
		return err
	}
	if parent.Err() != nil {
		return nil
	}
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	server := newHTTPServer(cfg, router, ctx)
	listener, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		logger.Error("startup_failed", "code", "LISTEN_FAILED")
		return errors.New("listener lokal tidak tersedia")
	}
	logger.Info("server_started", "environment", cfg.Environment, "addr", listener.Addr().String())
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		<-ctx.Done()
		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()
		if server.Shutdown(shutdownCtx) != nil {
			logger.Warn("shutdown_forced", "code", "SHUTDOWN_TIMEOUT")
			_ = server.Close()
		}
	}()
	err = server.Serve(listener)
	cancel()
	<-stopped
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server_failed", "code", "SERVER_FAILED")
		return errors.New("server lokal berhenti tidak normal")
	}
	logger.Info("server_stopped")
	return nil
}
