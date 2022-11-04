package realms

import (
	"context"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/domain/repositories"
	"github.com/alexZaicev/realm-mgr/internal/drivers/clock"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
	"github.com/alexZaicev/realm-mgr/internal/drivers/uuidgenerator"
)

type CreateRealmInput struct {
	Name        string
	Description string
}

func (i *CreateRealmInput) Validate() error {
	// TODO: add validation
	return nil
}

type CreateRealmRepos struct {
	Logger logging.Logger

	UUIDGen uuidgenerator.Generator
	Clock   clock.Clock

	Repository repositories.RealmManagerRepository
}

func (r *CreateRealmRepos) Validate() error {
	// TODO: add validation
	return nil
}

type CreateRealm struct {
}

func NewCreateRealm() *CreateRealm {
	return &CreateRealm{}
}

func (r *CreateRealm) CreateRealm(ctx context.Context, repos CreateRealmRepos, input CreateRealmInput) (entities.Realm, error) {
	if err := repos.Validate(); err != nil {
		return entities.Realm{}, nil
	}
	if err := input.Validate(); err != nil {
		return entities.Realm{}, nil
	}

	logger := repos.Logger.WithField("use-case", "create-realm")

	realmID, err := repos.UUIDGen.New()
	if err != nil {
		logger.WithError(err).Error("failed to generate UUID id")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to generate UUID id", nil)
	}

	now := repos.Clock.Now()

	realmToCreate := entities.Realm{
		ID:          realmID,
		Status:      entities.StatusDraft,
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if createErr := repos.Repository.CreateRealm(ctx, realmToCreate); createErr != nil {
		logger.WithError(err).Error("failed to create realm in repository")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to create realm in repository", nil)
	}

	return realmToCreate, nil
}
