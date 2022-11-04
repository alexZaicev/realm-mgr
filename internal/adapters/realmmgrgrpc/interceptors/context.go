package interceptors

import (
	"context"
	"reflect"

	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
)

type ctxKey string

const (
	loggerCtxKey ctxKey = "logger"
)

func LoggerFromContext(ctx context.Context) (logging.Logger, error) {
	logger, ok := ctx.Value(loggerCtxKey).(logging.Logger)
	if !ok || reflect.ValueOf(logger).IsNil() {
		return nil, realmmgr_errors.NewInternalError("logger not present in request context", nil)
	}

	return logger, nil
}

func ContextWithLogger(ctx context.Context, logger logging.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
