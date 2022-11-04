package realmmgrgrpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc/interceptors"
	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc/models"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
)

func (api *RealmManagerAPI) CreateRealm(
	ctx context.Context,
	req *realm_mgr_v1.CreateRealmRequest,
) (*realm_mgr_v1.CreateRealmResponse, error) {
	logger, err := interceptors.LoggerFromContext(ctx)
	if err != nil {
		api.backupLogger.WithError(err).Error("failed to extract logger from context")
		return nil, status.Errorf(codes.Internal, models.InternalErrMsg)
	}

	realm, err := api.realmOps.CreateRealm(ctx, logger, req.Name, req.Description)
	if err != nil {
		switch err.(type) {
		case *realmmgr_errors.InvalidArgumentError:
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, models.InternalErrMsg)
		}
	}

	grpcRealm, err := models.RealmFromDomain(realm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, models.InternalErrMsg)
	}

	return &realm_mgr_v1.CreateRealmResponse{
		Realm: grpcRealm,
	}, nil
}
