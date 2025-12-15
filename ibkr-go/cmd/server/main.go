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
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/api"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/database"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/session"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/telemetry"
	marketdatav1connect "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/marketdata/v1/marketdatav1connect"
	"github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/order/v1/orderv1connect"
	"github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/portfolio/v1/portfoliov1connect"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	shutdownTimeout     = 10 * time.Second
	otelShutdownTimeout = 5 * time.Second
	sessionTTL          = 24 * time.Hour
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

	// Initialize database.
	db, err := database.New(ctx, cfg.DBWriteDSN, cfg.DBReadDSN)
	if err != nil {
		logger.Error("Failed to initialize database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Database initialized successfully")

	// Initialize IBKR client.
	ibkrClient := ibkr.NewClient(cfg.IBKRGatewayURL, cfg.IBKRAccountID)

	// Initialize session service (24 hour TTL).
	sessionService := session.NewService(db, cfg.EncryptionKey, sessionTTL)

	logger.Info("Services initialized successfully")

	// Create and start HTTP server.
	server := setupServer(cfg, db, ibkrClient, sessionService)
	startServer(server, logger, cfg.MTLSEnabled)

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

func setupServer(
	cfg *config.Config,
	db *database.DB,
	ibkrClient *ibkr.Client,
	sessionService *session.Service,
) *http.Server {
	logger := slog.Default()
	mux := http.NewServeMux()

	// Create ConnectRPC interceptors.
	var interceptors connect.HandlerOption
	if cfg.MTLSEnabled {
		// Use mTLS authentication.
		interceptors = connect.WithInterceptors(
			middleware.NewValidationInterceptor(), // Validate requests with protovalidate.
			middleware.NewMTLSInterceptor(logger), // Validate mTLS client certificates.
			middleware.LoggingInterceptor(logger), // Log requests.
		)
	} else {
		// No authentication (development mode).
		interceptors = connect.WithInterceptors(
			middleware.NewValidationInterceptor(), // Validate requests with protovalidate.
			middleware.LoggingInterceptor(logger), // Log requests.
		)
	}

	// Create service handlers.
	orderHandler := api.NewOrderServiceHandler(ibkrClient)
	portfolioHandler := api.NewPortfolioServiceHandler(ibkrClient)
	marketDataHandler := api.NewMarketDataServiceHandler(ibkrClient)

	// Register service handlers.
	path, handler := orderv1connect.NewOrderServiceHandler(orderHandler, interceptors)
	mux.Handle(path, handler)

	path, handler = portfoliov1connect.NewPortfolioServiceHandler(portfolioHandler, interceptors)
	mux.Handle(path, handler)

	path, handler = marketdatav1connect.NewMarketDataServiceHandler(marketDataHandler, interceptors)
	mux.Handle(path, handler)

	logger.Info("Service handlers registered")

	// Health check endpoint.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Readiness check endpoint.
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		// Check database health.
		if err := db.Health(r.Context()); err != nil {
			logger.Error("Database health check failed", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "Database unhealthy")

			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ready")
	})

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)

	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// Configure TLS if enabled.
	if cfg.MTLSEnabled {
		if err := configureTLS(server, cfg); err != nil {
			logger.Error("Failed to configure TLS", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}

	return server
}

func startServer(server *http.Server, logger *slog.Logger, tlsEnabled bool) {
	go func() {
		if tlsEnabled {
			logger.Info("Server starting with mTLS", slog.String("addr", server.Addr))
			if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				logger.Error("Server failed", slog.String("error", err.Error()))
				os.Exit(1)
			}
		} else {
			logger.Info("Server starting", slog.String("addr", server.Addr))
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("Server failed", slog.String("error", err.Error()))
				os.Exit(1)
			}
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
