package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
)

type validatableAll interface {
	// ValidateAll is the newer method for validating gRPC requests and responses, added in
	// protoc-gen-validate v0.6.2.
	ValidateAll() error
}

type validatable interface {
	// Validate is the older method for validating gRPC requests and responses.
	Validate() error
}

// ValidateUnaryServerInterceptor validates the incoming request to a handler and the outgoing
// response from the handler. Validatable requests and responses can be created by importing
// "validate.proto" provided by https://github.com/envoyproxy/protoc-gen-validate. If a request or
// response does not have a method implemented for validation, the interceptor skips any validation
// attempts.
//
// If a request fails validation, an INVALID_ARGUMENT status code is returned along with the error
// message produced by the validation implementation. If the response fails validation, an INTERNAL
// status code is returned along with generic message indicating that an unexpected error occurred.
// In either failure scenario, a logger from the context is used to log the error so
// LoggerUnaryServerInterceptor must be executed before this interceptor.
func ValidateUnaryServerInterceptor(backupLogger logging.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger, err := LoggerFromContext(ctx)
		if err != nil {
			backupLogger.WithError(err).Error("failed to extract logger from context")
			return nil, status.Error(codes.Internal, "an unexpected error occurred")
		}

		if validateErr := validate(req); validateErr != nil {
			logger.WithError(validateErr).Warn("request validation failed")
			return nil, status.Error(codes.InvalidArgument, validateErr.Error())
		}

		// if the handler returns an error, then just propagate that back to the client
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		// if we fail to validate the response, then we have a bug
		// and we shouldn't let clients trip over themselves handling it
		if validateErr := validate(resp); validateErr != nil {
			logger.WithError(validateErr).Error("response validation failed")
			return nil, status.Error(codes.Internal, "an unexpected error occurred")
		}

		return resp, nil
	}
}

// validate attempts to validate the given input using one of the available validation interfaces.
// If the input does not implement either interface, no validation is performed.
func validate(input interface{}) error {
	if castInput, ok := input.(validatableAll); ok {
		return castInput.ValidateAll()
	}

	if castInput, ok := input.(validatable); ok {
		return castInput.Validate()
	}

	return nil
}
