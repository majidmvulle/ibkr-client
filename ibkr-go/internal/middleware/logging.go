package middleware

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
)

// LoggingInterceptor logs all incoming requests
func LoggingInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()

			logger.InfoContext(ctx, "request started",
				slog.String("procedure", req.Spec().Procedure),
				slog.String("protocol", req.Peer().Protocol),
			)

			resp, err := next(ctx, req)

			duration := time.Since(start)
			if err != nil {
				logger.ErrorContext(ctx, "request failed",
					slog.String("procedure", req.Spec().Procedure),
					slog.Duration("duration", duration),
					slog.String("error", err.Error()),
				)
			} else {
				logger.InfoContext(ctx, "request completed",
					slog.String("procedure", req.Spec().Procedure),
					slog.Duration("duration", duration),
				)
			}

			return resp, err
		}
	}
}
