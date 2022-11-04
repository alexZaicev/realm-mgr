package realms

import (
	"context"

	"github.com/google/uuid"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/domain/repositories"
	"github.com/alexZaicev/realm-mgr/internal/drivers/clock"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
)

type ReleaseRealmInput struct {
	RealmID uuid.UUID
}

func (i *ReleaseRealmInput) Validate() error {
	// TODO: add validation
	return nil
}

type ReleaseRealmRepos struct {
	Logger logging.Logger

	Clock clock.Clock

	Repository repositories.RealmManagerRepository
}

func (r *ReleaseRealmRepos) Validate() error {
	// TODO: add validation
	return nil
}

type ReleaseRealm struct {
}

func NewReleaseRealm() *ReleaseRealm {
	return &ReleaseRealm{}
}

func (r *ReleaseRealm) ReleaseRealm(ctx context.Context, repos ReleaseRealmRepos, input ReleaseRealmInput) (entities.Realm, error) {
	if err := repos.Validate(); err != nil {
		return entities.Realm{}, nil
	}
	if err := input.Validate(); err != nil {
		return entities.Realm{}, nil
	}

	logger := repos.Logger.WithFields(map[string]interface{}{
		"use-case": "release-realm",
		"realm-id": input.RealmID,
	})

	now := repos.Clock.Now()

	draftRealm, err := repos.Repository.GetRealm(ctx, input.RealmID, entities.StatusDraft)
	if err != nil {
		switch err.(type) {
		case *realmmgr_errors.NotFoundError:
			logger.WithError(err).Warn("no releasable realm found with provided ID")
			return entities.Realm{}, realmmgr_errors.NewNotFoundError("no releasable realm found with provided ID", nil)
		default:
			logger.WithError(err).Error("failed to get draft realm from repository")
			return entities.Realm{}, realmmgr_errors.NewInternalError("failed to get draft realm from repository", nil)
		}
	}

	activeRealm, err := repos.Repository.GetRealm(ctx, input.RealmID, entities.StatusActive)
	if err != nil {
		switch err.(type) {
		case *realmmgr_errors.NotFoundError:
			// it's a newly created realm that required update in status
			draftRealm.UpdatedAt = now
			draftRealm.Status = entities.StatusActive
			if updateErr := repos.Repository.UpdateRealm(ctx, draftRealm, entities.StatusDraft); updateErr != nil {
				logger.WithError(updateErr).Error("failed to update draft realm in repository")
				return entities.Realm{}, realmmgr_errors.NewInternalError("failed to update draft realm in repository", nil)
			}

			// TODO: perform other realm initializations

			return draftRealm, nil
		default:
			logger.WithError(err).Error("failed to get draft realm from repository")
			return entities.Realm{}, realmmgr_errors.NewInternalError("failed to get draft realm from repository", nil)
		}
	}

	activeRealm = activeRealm.Merge(draftRealm)
	activeRealm.UpdatedAt = now

	if deleteErr := repos.Repository.DeleteRealm(ctx, draftRealm.ID, draftRealm.Status); deleteErr != nil {
		logger.WithError(deleteErr).Error("failed to delete draft realm from repository")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to delete draft realm from repository", nil)
	}

	if updateErr := repos.Repository.UpdateRealm(ctx, activeRealm, activeRealm.Status); updateErr != nil {
		logger.WithError(updateErr).Error("failed to update active realm in repository")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to update active realm in repository", nil)
	}

	return activeRealm, nil
}
