package realmmgrgrpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/google/uuid"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
)

type RealmOps interface {
	GetRealm(ctx context.Context, logger logging.Logger, realmID uuid.UUID, status entities.Status) (entities.Realm, error)
	CreateRealm(ctx context.Context, logger logging.Logger, name, description string) (entities.Realm, error)
	ReleaseRealm(ctx context.Context, logger logging.Logger, realmID uuid.UUID) (entities.Realm, error)
	UpdateRealm(ctx context.Context, logger logging.Logger, realm entities.Realm) (entities.Realm, error)
}

type RealmManagerAPI struct {
	realm_mgr_v1.UnimplementedRealmManagerServiceServer

	backupLogger logging.Logger

	realmOps RealmOps
}

func NewRealmManagerAPI(
	backupLogger logging.Logger,
	realmOps RealmOps,
) (*RealmManagerAPI, error) {
	if backupLogger == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("backupLogger", realmmgr_errors.ErrMsgCannotBeNil)
	}
	if realmOps == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("realmOps", realmmgr_errors.ErrMsgCannotBeNil)
	}
	return &RealmManagerAPI{
		backupLogger: backupLogger,
		realmOps:     realmOps,
	}, nil
}

func (api *RealmManagerAPI) Register(server *grpc.Server) {
	realm_mgr_v1.RegisterRealmManagerServiceServer(server, api)
}
