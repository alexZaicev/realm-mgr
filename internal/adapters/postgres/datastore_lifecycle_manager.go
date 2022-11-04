package postgres

import (
	"context"

	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	"github.com/alexZaicev/realm-mgr/internal/drivers/pgdb"
	"github.com/alexZaicev/realm-mgr/internal/drivers/uuidgenerator"
)

type WriteConnProvider interface {
	TransactionalConn(ctx context.Context) (pgdb.Conn, error)
	Commit(conn pgdb.Conn) error
	Rollback(conn pgdb.Conn) error
}

type ReadConnProvider interface {
	Conn() pgdb.Conn
	TransactionalConn(ctx context.Context) (pgdb.Conn, error)
	Commit(conn pgdb.Conn) error
	Rollback(conn pgdb.Conn) error
}

type DataStoreLifecycleManager struct {
	readConnProvider  ReadConnProvider
	writeConnProvider WriteConnProvider
	uuidgen           uuidgenerator.Generator
}

func NewDataStoreLifecycleManager(
	readConnProvider ReadConnProvider,
	writeConnProvider WriteConnProvider,
	uuidgen uuidgenerator.Generator,
) (*DataStoreLifecycleManager, error) {
	if readConnProvider == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("readConnProvider", realmmgr_errors.ErrMsgCannotBeNil)
	}

	if writeConnProvider == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("writeConnProvider", realmmgr_errors.ErrMsgCannotBeNil)
	}

	if uuidgen == nil {
		return nil, realmmgr_errors.NewInvalidArgumentError("uuidgen", realmmgr_errors.ErrMsgCannotBeNil)
	}

	return &DataStoreLifecycleManager{
		readConnProvider:  readConnProvider,
		writeConnProvider: writeConnProvider,
		uuidgen:           uuidgen,
	}, nil
}

func (d *DataStoreLifecycleManager) NewNonTransactionalReadDatastore(_ context.Context) *DataStore {
	return &DataStore{
		uuidgen: d.uuidgen,
		db:      d.readConnProvider.Conn(),
	}
}

func (d *DataStoreLifecycleManager) NewWriteDatastore(ctx context.Context) (*DataStore, error) {
	tx, err := d.writeConnProvider.TransactionalConn(ctx)
	if err != nil {
		return nil, err
	}
	return &DataStore{
		uuidgen: d.uuidgen,
		db:      tx,
	}, nil
}

func (d *DataStoreLifecycleManager) CommitChanges(datastore *DataStore) error {
	return d.writeConnProvider.Commit(datastore.db)
}

func (d *DataStoreLifecycleManager) RollbackChanges(datastore *DataStore) error {
	return d.writeConnProvider.Rollback(datastore.db)
}
