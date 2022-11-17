package realmmgrgrpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc/interceptors"
	"github.com/alexZaicev/realm-mgr/internal/adapters/realmmgrgrpc/models"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
)

func (api *RealmManagerAPI) ReleaseRealm(
	ctx context.Context,
	req *realm_mgr_v1.ReleaseRealmRequest,
) (*realm_mgr_v1.ReleaseRealmResponse, error) {
	logger, err := interceptors.LoggerFromContext(ctx)
	if err != nil {
		api.backupLogger.WithError(err).Error("failed to extract logger from context")
		return nil, status.Errorf(codes.Internal, models.InternalErrMsg)
	}

	realmID, err := uuid.Parse(req.Id)
	if err != nil {
		logger.WithError(err).WithField("realm-id", req.Id).Info("invalid realm ID supplied")
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf(models.InvalidRealmID, req.Id))
	}

	realm, err := api.realmOps.ReleaseRealm(ctx, logger, realmID)
	if err != nil {
		switch err.(type) {
		case *realmmgr_errors.NotFoundError:
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("no releasable realm with ID found: %s", realmID))
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

	return &realm_mgr_v1.ReleaseRealmResponse{
		Realm: grpcRealm,
	}, nil
}
