package realmmgrgrpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc/interceptors"
	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc/models"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
)

func (api *RealmManagerAPI) UpdateRealm(
	ctx context.Context,
	req *realm_mgr_v1.UpdateRealmRequest,
) (*realm_mgr_v1.UpdateRealmResponse, error) {
	logger, err := interceptors.LoggerFromContext(ctx)
	if err != nil {
		api.backupLogger.WithError(err).Error("failed to extract logger from context")
		return nil, status.Errorf(codes.Internal, models.InternalErrMsg)
	}

	realmInput, err := models.RealmToDomain(req.Realm)
	if err != nil {
		logger.WithError(err).Info("invalid realm data supplied")
		return nil, status.Errorf(codes.InvalidArgument, "invalid realm data supplied")
	}

	realm, err := api.realmOps.UpdateRealm(ctx, logger, realmInput)
	if err != nil {
		switch err.(type) {
		case *realmmgr_errors.NotFoundError:
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("realm with ID not found: %s", realmInput.ID))
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

	return &realm_mgr_v1.UpdateRealmResponse{
		Realm: grpcRealm,
	}, nil
}
