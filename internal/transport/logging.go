package transport

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func unaryLoggingInterceptor(logger *slog.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		logger.Debug("grpc unary call",
			"method", method,
			"duration", time.Since(start),
			"error", errString(err),
		)
		return err
	}
}

func streamLoggingInterceptor(logger *slog.Logger) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		start := time.Now()
		cs, err := streamer(ctx, desc, cc, method, opts...)
		logger.Debug("grpc stream call",
			"method", method,
			"duration", time.Since(start),
			"error", errString(err),
		)
		return cs, err
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
