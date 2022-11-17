package realms

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/domain/repositories"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
)

type GetRealmInput struct {
	RealmID uuid.UUID
	Status  entities.Status
}

func (i *GetRealmInput) Validate() error {
	// TODO: add validation
	return nil
}

type GetRealmRepos struct {
	Logger logging.Logger

	Repository repositories.RealmManagerRepository
}

func (r *GetRealmRepos) Validate() error {
	// TODO: add validation
	return nil
}

type GetRealm struct {
}

func NewGetRealm() *GetRealm {
	return &GetRealm{}
}

func (r *GetRealm) GetRealm(ctx context.Context, repos GetRealmRepos, input GetRealmInput) (entities.Realm, error) {
	if err := repos.Validate(); err != nil {
		return entities.Realm{}, nil
	}
	if err := input.Validate(); err != nil {
		return entities.Realm{}, nil
	}

	logger := repos.Logger.WithFields(map[string]interface{}{
		"use-case": "get-realm",
		"realm-id": input.RealmID,
	})

	realm, err := repos.Repository.GetRealm(ctx, input.RealmID, input.Status)
	if err != nil {
		switch err.(type) {
		case *realmmgr_errors.NotFoundError:
			return entities.Realm{}, realmmgr_errors.NewNotFoundError(
				fmt.Sprintf("realm with ID %s not found", input.RealmID),
				nil,
			)
		default:
			logger.WithError(err).Error("failed to get realm from repository")
			return entities.Realm{}, realmmgr_errors.NewInternalError("failed to get realm from repository", nil)
		}
	}

	// return NotFoundError in case realm is marked as deleted
	if realm.Status == entities.StatusDeleted {
		return entities.Realm{}, realmmgr_errors.NewNotFoundError(
			fmt.Sprintf("realm with ID %s not found", input.RealmID),
			nil,
		)
	}

	return realm, nil
}
