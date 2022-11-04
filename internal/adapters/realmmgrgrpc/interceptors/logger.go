package interceptors

import (
	"context"
	"path"

	"google.golang.org/grpc"

	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
)

const (
	serviceKey = "grpc-service"
	methodKey  = "grpc-method"
)

// LoggerUnaryServerInterceptor creates a new logger instance and adds it to
// the context such that it can be accessed by other interceptors and the
// request handler.
//
// Each request received by the interceptor will be given a logger with a
// service and method field applied to easily identify the gRPC API that was
// called via structured logging.
func LoggerUnaryServerInterceptor(logger logging.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		service := path.Dir(info.FullMethod)[1:]
		method := path.Base(info.FullMethod)

		requestLogger := logger.WithFields(logging.Fields{
			serviceKey: service,
			methodKey:  method,
		})

		ctx = ContextWithLogger(ctx, requestLogger)
		return handler(ctx, req)
	}
}
