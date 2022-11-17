package utils

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/alexZaicev/realm-mgr/internal/drivers/config"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
)

const (
	ConfigGRPCTarget       = "grpc.target"
	DefaultGRPCCallTimeout = time.Second
)

func NewRealmManagerGRPCClient(cfg config.Config) (realm_mgr_v1.RealmManagerServiceClient, error) {
	target, err := config.Get[string](cfg, ConfigGRPCTarget)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultGRPCCallTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithReturnConnectionError(),
	}...)
	if err != nil {
		return nil, err
	}
	return realm_mgr_v1.NewRealmManagerServiceClient(conn), nil
}

func MakeGRPCRequestContext(ctx context.Context, headers ...string) (context.Context, error) {
	if len(headers)%2 == 1 {
		return nil, fmt.Errorf("headers should have even element count")
	}
	return metadata.AppendToOutgoingContext(ctx), nil
}
