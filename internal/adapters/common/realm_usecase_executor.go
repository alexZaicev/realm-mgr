package common

import (
	"context"

	"github.com/google/uuid"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/drivers/clock"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
	"github.com/alexZaicev/realm-mgr/internal/drivers/uuidgenerator"
	"github.com/alexZaicev/realm-mgr/internal/usecases/realms"
)

type RealmGetter interface {
	GetRealm(ctx context.Context, repos realms.GetRealmRepos, input realms.GetRealmInput) (entities.Realm, error)
}

type RealmCreator interface {
	CreateRealm(ctx context.Context, repos realms.CreateRealmRepos, input realms.CreateRealmInput) (entities.Realm, error)
}

type RealmReleaser interface {
	ReleaseRealm(ctx context.Context, repos realms.ReleaseRealmRepos, input realms.ReleaseRealmInput) (entities.Realm, error)
}

type RealmUseCaseExecutor struct {
	uuidGen          uuidgenerator.Generator
	clock            clock.Clock
	dataStoreManager DataStoreManager

	realmGetter   RealmGetter
	realmCreator  RealmCreator
	realmReleaser RealmReleaser
}

func NewRealmUseCaseExecutor(
	uuidGen uuidgenerator.Generator,
	clock clock.Clock,
	dataStoreManager DataStoreManager,
	realmGetter RealmGetter,
	realmCreator RealmCreator,
	realmReleaser RealmReleaser,
) (*RealmUseCaseExecutor, error) {
	if uuidGen == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("uuidGen", realmmgr_errors.ErrMsgCannotBeNil)
	}
	if clock == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("clock", realmmgr_errors.ErrMsgCannotBeNil)
	}
	if dataStoreManager == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("dataStoreManager", realmmgr_errors.ErrMsgCannotBeNil)
	}
	if realmGetter == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("realmGetter", realmmgr_errors.ErrMsgCannotBeNil)
	}
	if realmCreator == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("realmCreator", realmmgr_errors.ErrMsgCannotBeNil)
	}
	if realmReleaser == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("realmReleaser", realmmgr_errors.ErrMsgCannotBeNil)
	}
	return &RealmUseCaseExecutor{
		uuidGen:          uuidGen,
		clock:            clock,
		dataStoreManager: dataStoreManager,
		realmGetter:      realmGetter,
		realmCreator:     realmCreator,
		realmReleaser:    realmReleaser,
	}, nil
}

func (e *RealmUseCaseExecutor) GetRealm(ctx context.Context, logger logging.Logger, realmID uuid.UUID, status entities.Status) (entities.Realm, error) {
	repository := e.dataStoreManager.NewNonTransactionalReadDatastore(ctx)

	repos := realms.GetRealmRepos{
		Logger:     logger,
		Repository: repository,
	}

	input := realms.GetRealmInput{
		RealmID: realmID,
		Status:  status,
	}

	realm, err := e.realmGetter.GetRealm(ctx, repos, input)
	if err != nil {
		return entities.Realm{}, err
	}

	return realm, nil
}

func (e *RealmUseCaseExecutor) CreateRealm(ctx context.Context, logger logging.Logger, name, description string) (entities.Realm, error) {
	repository, auditRepository, rollbackFn, err := TransactionalRepositories(ctx, logger, e.dataStoreManager)
	if err != nil {
		logger.WithError(err).Error("failed to configure repositories")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to configure repositories", err)
	}
	defer rollbackFn()

	repos := realms.CreateRealmRepos{
		Logger:     logger,
		UUIDGen:    e.uuidGen,
		Clock:      e.clock,
		Repository: repository,
	}

	input := realms.CreateRealmInput{
		Name:        name,
		Description: description,
	}

	realm, err := e.realmCreator.CreateRealm(ctx, repos, input)
	if err != nil {
		return entities.Realm{}, err
	}

	if commitErr := CommitRepositories(logger, e.dataStoreManager, repository, auditRepository); err != nil {
		logger.WithError(commitErr).Error("failed to commit transaction")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to commit transaction", commitErr)
	}

	return realm, nil
}

func (e *RealmUseCaseExecutor) ReleaseRealm(ctx context.Context, logger logging.Logger, realmID uuid.UUID) (entities.Realm, error) {
	repository, auditRepository, rollbackFn, err := TransactionalRepositories(ctx, logger, e.dataStoreManager)
	if err != nil {
		logger.WithError(err).Error("failed to configure repositories")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to configure repositories", err)
	}
	defer rollbackFn()

	repos := realms.ReleaseRealmRepos{
		Logger:     logger,
		Clock:      e.clock,
		Repository: repository,
	}

	input := realms.ReleaseRealmInput{
		RealmID: realmID,
	}

	realm, err := e.realmReleaser.ReleaseRealm(ctx, repos, input)
	if err != nil {
		return entities.Realm{}, err
	}

	if commitErr := CommitRepositories(logger, e.dataStoreManager, repository, auditRepository); err != nil {
		logger.WithError(commitErr).Error("failed to commit transaction")
		return entities.Realm{}, realmmgr_errors.NewInternalError("failed to commit transaction", commitErr)
	}

	return realm, nil
}
