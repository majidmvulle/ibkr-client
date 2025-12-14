package middleware

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

// NewValidationInterceptor creates a ConnectRPC interceptor that validates requests using protovalidate.
func NewValidationInterceptor() connect.UnaryInterceptorFunc {
	validator, err := protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize validator: %v", err))
	}

	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Validate the request message.
			if msg, ok := req.Any().(proto.Message); ok {
				if err := validator.Validate(msg); err != nil {
					return nil, connect.NewError(connect.CodeInvalidArgument, err)
				}
			}

			return next(ctx, req)
		}
	}
}
