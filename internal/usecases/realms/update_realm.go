package realms

import (
	"context"
	"fmt"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/domain/repositories"
	realmmgr_clock "github.com/alexZaicev/realm-mgr/internal/drivers/clock"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
)

type UpdateRealmInput struct {
	Realm entities.Realm
}

func (i *UpdateRealmInput) Validate() error {
	// TODO: add validation
	return nil
}

type UpdateRealmRepos struct {
	Logger     logging.Logger
	Clock      realmmgr_clock.Clock
	Repository repositories.RealmManagerRepository
}

func (r *UpdateRealmRepos) Validate() error {
	// TODO: add validation
	return nil
}

type UpdateRealm struct {
}

func NewUpdateRealm() *UpdateRealm {
	return &UpdateRealm{}
}

func (r *UpdateRealm) UpdateRealm(ctx context.Context, repos UpdateRealmRepos, input UpdateRealmInput) (entities.Realm, error) {
	if err := repos.Validate(); err != nil {
		return entities.Realm{}, nil
	}
	if err := input.Validate(); err != nil {
		return entities.Realm{}, nil
	}

	logger := repos.Logger.WithFields(map[string]interface{}{
		"use-case": "update-realm",
		"realm-id": input.Realm.ID,
	})

	now := repos.Clock.Now()
	input.Realm.UpdatedAt = now

	// check if draft for realm already exists
	draftRealm, err := repos.Repository.GetRealm(ctx, input.Realm.ID, entities.StatusDraft)
	if err != nil {
		switch err.(type) {
		case *realmmgr_errors.NotFoundError:
			// create a new draft
			return r.createNewDraft(ctx, logger, repos, input)
		default:
			logger.WithError(err).Error("failed to get draft realm from repository")
			return entities.Realm{}, realmmgr_errors.NewInternalError("failed to get draft realm from repository", nil)
		}
	}

	// update existing draft
	draftRealm = draftRealm.Merge(input.Realm)

	if updateErr := repos.Repository.UpdateRealm(ctx, draftRealm, draftRealm.Status); updateErr != nil {
		logger.WithError(err).Error("failed to update realm in repository")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to update realm in repository", nil)
	}

	return draftRealm, nil
}

func (r *UpdateRealm) createNewDraft(
	ctx context.Context,
	logger logging.Logger,
	repos UpdateRealmRepos,
	input UpdateRealmInput,
) (entities.Realm, error) {
	activeRealm, err := repos.Repository.GetRealm(ctx, input.Realm.ID, entities.StatusActive)
	if err != nil {
		switch err.(type) {
		case *realmmgr_errors.NotFoundError:
			return entities.Realm{}, realmmgr_errors.NewNotFoundError(
				fmt.Sprintf("realm with ID %s not found", input.Realm.ID),
				nil,
			)
		default:
			logger.WithError(err).Error("failed to get active realm from repository")
			return entities.Realm{}, realmmgr_errors.NewInternalError("failed to get active realm from repository", nil)
		}
	}

	draftRealm := activeRealm.Merge(input.Realm)
	draftRealm.Status = entities.StatusDraft

	if createErr := repos.Repository.CreateRealm(ctx, draftRealm); createErr != nil {
		logger.WithError(err).Error("failed to create draft realm in repository")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to create draft realm in repository", nil)
	}

	return draftRealm, nil
}
