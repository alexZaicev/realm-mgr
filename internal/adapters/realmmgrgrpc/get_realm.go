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

func (api *RealmManagerAPI) GetRealm(
	ctx context.Context,
	req *realm_mgr_v1.GetRealmRequest,
) (*realm_mgr_v1.GetRealmResponse, error) {
	logger, err := interceptors.LoggerFromContext(ctx)
	if err != nil {
		api.backupLogger.WithError(err).Error("failed to extract logger from context")
		return nil, status.Errorf(codes.Internal, models.InternalErrMsg)
	}

	realmID, err := uuid.Parse(req.Id)
	if err != nil {
		logger.WithError(err).WithField("realm-id", req.Id).Info("invalid realm ID supplied")
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("realm ID was not a valid UUID: %s", req.Id))
	}

	// if status not provided in the request, default it to always return active realm
	if req.Status == realm_mgr_v1.EnumStatus_ENUM_STATUS_UNSPECIFIED {
		req.Status = realm_mgr_v1.EnumStatus_ENUM_STATUS_ACTIVE
	}

	realmStatus, ok := models.StatusGRPCValues[req.Status]
	if !ok {
		logger.WithError(err).WithField("status", req.Status).Info("invalid realm status supplied")
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("unexpected realm status: %s", req.Status))
	}

	realm, err := api.realmOps.GetRealm(ctx, logger, realmID, realmStatus)
	if err != nil {
		switch err.(type) {
		case *realmmgr_errors.NotFoundError:
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("realm with ID not found: %s", req.Id))
		default:
			return nil, status.Errorf(codes.Internal, models.InternalErrMsg)
		}
	}

	grpcRealm, err := models.RealmFromDomain(realm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, models.InternalErrMsg)
	}

	return &realm_mgr_v1.GetRealmResponse{
		Realm: grpcRealm,
	}, nil
}
