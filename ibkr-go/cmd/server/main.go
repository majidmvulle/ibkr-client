package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/telemetry"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up structured logging
	logLevel := slog.LevelInfo
	if cfg.AppDebug {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting IBKR client",
		slog.String("version", cfg.AppVersion),
		slog.String("env", cfg.AppEnv),
	)

	// Initialize OpenTelemetry
	tp, err := telemetry.InitTracer(ctx, cfg.AppName, cfg.OtelCollectorEndpoint)
	if err != nil {
		logger.Error("Failed to initialize tracer", slog.String("error", err.Error()))
	} else {
		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := telemetry.Shutdown(shutdownCtx, tp); err != nil {
				logger.Error("Failed to shutdown tracer", slog.String("error", err.Error()))
			}
		}()
		logger.Info("OpenTelemetry initialized", slog.String("endpoint", cfg.OtelCollectorEndpoint))
	}

	// TODO: Initialize database connection
	// TODO: Initialize IBKR client
	// TODO: Set up ConnectRPC handlers

	// Create HTTP server with h2c (HTTP/2 without TLS for local dev)
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Readiness check endpoint
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Check database connection, IBKR gateway, etc.
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ready")
	})

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server starting", slog.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	logger.Info("Server stopped")
}
