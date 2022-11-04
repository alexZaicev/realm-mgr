package grpcserver

import (
	"context"

	"google.golang.org/grpc/metadata"

	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/drivers/headers"
)

// MetadataCarrier implements headers.Carrier for headers stored in a gRPC context as provided by
// https://pkg.go.dev/google.golang.org/grpc/metadata#MD.
type MetadataCarrier struct {
	metadata metadata.MD
}

// NewMetadataCarrierFromIncomingContext creates a new MetadataCarrier, using the gRPC metadata
// found in an incoming request context.
func NewMetadataCarrierFromIncomingContext(ctx context.Context) (*MetadataCarrier, error) {
	if ctx == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("ctx", "cannot be nil")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, realmmgr_errors.NewInvalidArgumentError("ctx", "does not contain metadata")
	}

	return &MetadataCarrier{metadata: md}, nil
}

// Metadata returns a copy of the wrapped gRPC metadata.
func (c *MetadataCarrier) Metadata() metadata.MD {
	return c.metadata.Copy()
}

// Non-allocating compile time check to ensure the headers.Carrier interface is implemented
// correctly.
var _ headers.Carrier = (*MetadataCarrier)(nil)

// GetSingle reads a single value from the requested header.
//
// If no value is found, a headers.HeaderNotFound error is returned. If multiple values are found, a
// headers.MultipleHeadersFound error is returned.
func (c *MetadataCarrier) GetSingle(key string) (string, error) {
	values := c.metadata.Get(key)
	if len(values) == 0 {
		return "", headers.NewHeaderNotFound(key)
	}

	if len(values) > 1 {
		return "", headers.NewMultipleHeadersFound(key)
	}

	return values[0], nil
}

// GetMultiple reads multiple values from the requested header.
func (c *MetadataCarrier) GetMultiple(key string) []string {
	values := c.metadata.Get(key)

	// normalise not found to empty
	if values == nil {
		return []string{}
	}

	return values
}

// Set sets the given value for a header, removing existing values in the
// process.
func (c *MetadataCarrier) Set(key, value string) {
	c.metadata.Set(key, value)
}
