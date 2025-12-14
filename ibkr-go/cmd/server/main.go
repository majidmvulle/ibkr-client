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

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/telemetry"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	shutdownTimeout     = 10 * time.Second
	otelShutdownTimeout = 5 * time.Second
)

func main() {
	ctx := context.Background()

	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up structured logging.
	logger := setupLogger(cfg)

	logger.Info("Starting IBKR client",
		slog.String("version", cfg.AppVersion),
		slog.String("env", cfg.AppEnv),
	)

	// Initialize OpenTelemetry.
	tp := initTelemetry(ctx, cfg, logger)
	defer shutdownTelemetry(tp, logger)

	// Create and start HTTP server.
	server := setupServer(cfg)
	startServer(server, logger)

	// Wait for shutdown signal.
	waitForShutdown(server, logger)
}

func setupLogger(cfg *config.Config) *slog.Logger {
	logLevel := slog.LevelInfo
	if cfg.AppDebug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	return logger
}

func initTelemetry(ctx context.Context, cfg *config.Config, logger *slog.Logger) *sdktrace.TracerProvider {
	tp, err := telemetry.InitTracer(ctx, cfg.AppName, cfg.OtelCollectorEndpoint)
	if err != nil {
		logger.Error("Failed to initialize tracer", slog.String("error", err.Error()))

		return nil
	}

	logger.Info("OpenTelemetry initialized", slog.String("endpoint", cfg.OtelCollectorEndpoint))

	return tp
}

func shutdownTelemetry(tp *sdktrace.TracerProvider, logger *slog.Logger) {
	if tp == nil {
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), otelShutdownTimeout)
	defer cancel()

	if err := telemetry.Shutdown(shutdownCtx, tp); err != nil {
		logger.Error("Failed to shutdown tracer", slog.String("error", err.Error()))
	}
}

func setupServer(cfg *config.Config) *http.Server {
	logger := slog.Default()
	mux := http.NewServeMux()

	// Create ConnectRPC interceptors.
	interceptors := connect.WithInterceptors(
		middleware.NewValidationInterceptor(), // Validate requests with protovalidate.
		middleware.LoggingInterceptor(logger), // Log requests.
	)

	// Service handlers will be registered here in Phase 4.
	// Example:
	// mux.Handle(orderv1connect.NewOrderServiceHandler(orderHandler, interceptors)).
	// mux.Handle(portfoliov1connect.NewPortfolioServiceHandler(portfolioHandler, interceptors)).
	// mux.Handle(marketdatav1connect.NewMarketDataServiceHandler(marketDataHandler, interceptors)).
	_ = interceptors // Will be used when service handlers are implemented..

	// Health check endpoint.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Readiness check endpoint.
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ready")
	})

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)

	return &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
}

func startServer(server *http.Server, logger *slog.Logger) {
	go func() {
		logger.Info("Server starting", slog.String("addr", server.Addr))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()
}

func waitForShutdown(server *http.Server, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	logger.Info("Server stopped")
}
