package grpcserver_test

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/drivers/grpcserver"
	"github.com/alexZaicev/realm-mgr/internal/drivers/headers"
)

func Test_NewMetadataCarrierFromIncomingContext_Happy(t *testing.T) {
	// arrange
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})

	// act
	carrier, err := grpcserver.NewMetadataCarrierFromIncomingContext(ctx)

	// assert
	assert.NotNil(t, carrier)
	assert.NoError(t, err)
}

func Test_NewMetadataCarrierFromIncomingContext_InvalidParameter(t *testing.T) {
	testCases := []struct {
		name   string
		ctx    context.Context
		errMsg string
	}{
		{
			name:   "nil ctx",
			ctx:    nil,
			errMsg: "invalid parameter ctx: cannot be nil",
		},
		{
			name:   "missing metadata",
			ctx:    context.Background(),
			errMsg: "invalid parameter ctx: does not contain metadata",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// no arrange necessary

			// act
			carrier, err := grpcserver.NewMetadataCarrierFromIncomingContext(tc.ctx)

			// assert
			assert.Nil(t, carrier)
			require.EqualError(t, err, tc.errMsg)
			assert.IsType(t, realmmgr_errors.InvalidArgumentErrorType, err)
		})
	}
}

func Test_MetadataCarrier_GetSingle_Happy(t *testing.T) {
	// arrange
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.MD{"example": []string{"value"}},
	)

	carrier, err := grpcserver.NewMetadataCarrierFromIncomingContext(ctx)
	require.NoError(t, err)

	// act
	value, err := carrier.GetSingle("example")

	// assert
	assert.Equal(t, "value", value)
	assert.NoError(t, err)
}

func Test_MetadataCarrier_GetSingle_Error(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     context.Context
		errMsg  string
		errType error
	}{
		{
			name: "header not found",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{},
			),
			errMsg:  "header was not found with name: example",
			errType: headers.HeaderNotFoundType,
		},
		{
			name: "multiple headers found",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{"example": []string{"value", "value 2"}},
			),
			errMsg:  `invalid input(s): header "example" cannot have multiple values`,
			errType: headers.MultipleHeadersFoundType,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			carrier, err := grpcserver.NewMetadataCarrierFromIncomingContext(tc.ctx)
			require.NoError(t, err)

			// act
			value, err := carrier.GetSingle("example")

			// assert
			assert.Empty(t, value)
			require.EqualError(t, err, tc.errMsg)
			assert.IsType(t, tc.errType, err)
		})
	}
}

func Test_MetadataCarrier_GetMultiple(t *testing.T) {
	testCases := []struct {
		name     string
		ctx      context.Context
		expected []string
	}{
		{
			name: "missing header",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{},
			),
			expected: []string{},
		},
		{
			name: "empty header",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{"example": []string{}},
			),
			expected: []string{},
		},
		{
			name: "single header value",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{"example": []string{"value"}},
			),
			expected: []string{"value"},
		},
		{
			name: "multiple header values",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{"example": []string{"value", "value 2", "value 3"}},
			),
			expected: []string{"value", "value 2", "value 3"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			carrier, err := grpcserver.NewMetadataCarrierFromIncomingContext(tc.ctx)
			require.NoError(t, err)

			// act
			values := carrier.GetMultiple("example")

			// assert
			assert.Equal(t, tc.expected, values)
		})
	}
}

func Test_MetadataCarrier_Set(t *testing.T) {
	testCases := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "unset header",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{},
			),
		},
		{
			name: "empty header",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{"example": []string{}},
			),
		},
		{
			name: "single header value",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{"example": []string{"value"}},
			),
		},
		{
			name: "multiple header values",
			ctx: metadata.NewIncomingContext(
				context.Background(),
				metadata.MD{"example": []string{"value", "value 2", "value 3"}},
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			carrier, err := grpcserver.NewMetadataCarrierFromIncomingContext(tc.ctx)
			require.NoError(t, err)

			// act
			carrier.Set("example", "other value")

			// assert
			md := carrier.Metadata()

			values := md.Get("example")
			assert.Equal(t, []string{"other value"}, values)
		})
	}
}
