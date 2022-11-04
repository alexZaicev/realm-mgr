package common

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alexZaicev/realm-mgr/internal/adapters/postgres"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/domain/repositories"
	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
)

type DataStoreManager interface {
	NewNonTransactionalReadDatastore(ctx context.Context) repositories.RealmManagerRepository
	NewNonTransactionalReadAuditstore(ctx context.Context) repositories.RealmManagerAuditRepository
	NewWriteDatastore(ctx context.Context) (repositories.RealmManagerRepository, repositories.RealmManagerAuditRepository, error)
	CommitChanges(datastore repositories.RealmManagerRepository, auditstore repositories.RealmManagerAuditRepository) error
	RollbackChanges(datastore repositories.RealmManagerRepository, auditstore repositories.RealmManagerAuditRepository) error
}

type PgDatastoreLifeCycleManager interface {
	NewNonTransactionalReadDatastore(context.Context) *postgres.DataStore
	NewWriteDatastore(context.Context) (*postgres.DataStore, error)
	CommitChanges(*postgres.DataStore) error
	RollbackChanges(*postgres.DataStore) error
}

type PgDataStoreManager struct {
	DatastoreLifeCycleManager PgDatastoreLifeCycleManager
}

func NewPgDataStoreManager(datastoreLifeCycleManager PgDatastoreLifeCycleManager) (*PgDataStoreManager, error) {
	if datastoreLifeCycleManager == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("datastoreLifeCycleManager", realmmgr_errors.ErrMsgCannotBeNil)
	}
	return &PgDataStoreManager{DatastoreLifeCycleManager: datastoreLifeCycleManager}, nil
}

func (pdm *PgDataStoreManager) NewNonTransactionalReadDatastore(ctx context.Context) repositories.RealmManagerRepository {
	return pdm.DatastoreLifeCycleManager.NewNonTransactionalReadDatastore(ctx)
}

func (pdm *PgDataStoreManager) NewNonTransactionalReadAuditstore(ctx context.Context) repositories.RealmManagerAuditRepository {
	return pdm.DatastoreLifeCycleManager.NewNonTransactionalReadDatastore(ctx)
}

func (pdm *PgDataStoreManager) NewWriteDatastore(
	ctx context.Context,
) (repositories.RealmManagerRepository, repositories.RealmManagerAuditRepository, error) {
	pgstore, err := pdm.DatastoreLifeCycleManager.NewWriteDatastore(ctx)
	if err != nil {
		return nil, nil, err
	}
	return pgstore, pgstore, nil
}

func (pdm *PgDataStoreManager) CommitChanges(
	datastore repositories.RealmManagerRepository,
	auditstore repositories.RealmManagerAuditRepository,
) error {
	castDatastore, ok := datastore.(*postgres.DataStore)
	if !ok {
		return realmmgr_errors.NewInvalidArgumentError("datastore", "underlying type must be *postgres.DataStore")
	}

	castAuditstore, ok := auditstore.(*postgres.DataStore)
	if !ok {
		return realmmgr_errors.NewInvalidArgumentError("auditstore", "underlying type must be *postgres.DataStore")
	}
	if castDatastore != castAuditstore {
		return realmmgr_errors.NewInvalidArgumentError("auditstore", "must reference the same underlying object as datastore")
	}
	return pdm.DatastoreLifeCycleManager.CommitChanges(castDatastore)
}

func (pdm *PgDataStoreManager) RollbackChanges(
	datastore repositories.RealmManagerRepository,
	auditstore repositories.RealmManagerAuditRepository,
) error {
	castDatastore, ok := datastore.(*postgres.DataStore)
	if !ok {
		return realmmgr_errors.NewInvalidArgumentError("datastore", "underlying type must be *postgres.DataStore")
	}

	castAuditstore, ok := auditstore.(*postgres.DataStore)
	if !ok {
		return realmmgr_errors.NewInvalidArgumentError("auditstore", "underlying type must be *postgres.DataStore")
	}
	if castDatastore != castAuditstore {
		return realmmgr_errors.NewInvalidArgumentError("auditstore", "must reference the same underlying object as datastore")
	}
	return pdm.DatastoreLifeCycleManager.RollbackChanges(castDatastore)
}

const RollbackTimeoutSeconds = 30

func TransactionalRepositories(
	ctx context.Context,
	logger logging.Logger,
	dataStoreManager DataStoreManager,
) (
	repositories.RealmManagerRepository,
	repositories.RealmManagerAuditRepository,
	func(),
	error,
) {
	datastoreTx, auditstoreTx, err := dataStoreManager.NewWriteDatastore(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to configure a new datastore and auditstore")
		return nil, nil, nil, err
	}

	rollbackFunc := func() {
		if rollbackErr := dataStoreManager.RollbackChanges(datastoreTx, auditstoreTx); rollbackErr != nil {
			if errors.Is(rollbackErr, sql.ErrTxDone) {
				// this is an expected error and can be ignored
				return
			}
			logger.WithError(rollbackErr).Error("Failed to rollback db changes")
		}
	}

	return datastoreTx, auditstoreTx, rollbackFunc, nil
}

func CommitRepositories(
	logger logging.Logger,
	dataStoreManager DataStoreManager,
	datastore repositories.RealmManagerRepository,
	auditstore repositories.RealmManagerAuditRepository,
) error {
	if err := dataStoreManager.CommitChanges(datastore, auditstore); err != nil {
		logger.WithError(err).Error("Failed to commit database transaction")
		return err
	}

	return nil
}
