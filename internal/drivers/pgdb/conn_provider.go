package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// but is needed by the database/sql package
	_ "github.com/jackc/pgx/v5/stdlib"

	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
)

// Conn provides an abstraction over the sql.DB https://pkg.go.dev/database/sql#DB and
// sql.Tx https://pkg.go.dev/database/sql#Tx types to allow for their interchangeable use within implemented functionality.
// This interface contains both the context and no context version of operations as some tools
// expect the no context version of the functions. However, use of the context versions
// wherever possible is strongly encouraged.
type Conn interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// ConnProvider provides transactional and non-transactional connections on demand.
type ConnProvider struct {
	sqlDB *sql.DB
	hooks []Hook
}

// NewConnProvider constructs a new provider for the database specified in the config
// parameters using the pgx as the underlying driver.
// NOTE: SSL mode should never be disabled in production.
func NewConnProvider(
	host,
	port,
	user,
	password,
	databaseName,
	sslMode string,
	opts ...DBOption,
) (*ConnProvider, error) {
	openCmd := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		databaseName,
		sslMode,
	)

	sqlDB, err := sql.Open("pgx", openCmd)
	if err != nil {
		// Not possible to reach in unit tests
		return nil, realmmgr_errors.NewInternalError("unable to connect to database", err)
	}

	// db options that are not configured are left as nil pointers to ensure
	// that we don't modify the default behavior of `sql.DB` unless specified by
	// the caller.
	options := &DBOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}

	// if an options pointer is not nil, it has been configured and is safe to
	// dereference to modify the underlying `sql.DB` instance.
	// set a default timeout to protect us from idle connections being closed by Istio
	// see the comments on `defaultMaxIdleTimeMinutes` for more details
	sqlDB.SetConnMaxIdleTime(time.Minute * time.Duration(defaultMaxIdleTimeMinutes))
	if options.MaxIdleTime != nil {
		sqlDB.SetConnMaxIdleTime(*options.MaxIdleTime)
	}

	if options.MaxLifetime != nil {
		sqlDB.SetConnMaxLifetime(*options.MaxLifetime)
	}

	if options.MaxOpenConns != nil {
		sqlDB.SetMaxOpenConns(*options.MaxOpenConns)
	}

	if options.MaxIdleConns != nil {
		sqlDB.SetMaxIdleConns(*options.MaxIdleConns)
	}

	return &ConnProvider{
		sqlDB: sqlDB,
		hooks: options.Hooks,
	}, nil
}

// Conn returns a non-transactional connection to be used to query the database.
func (d *ConnProvider) Conn() Conn {
	return &hookableConn{
		db:    d.sqlDB,
		hooks: d.hooks,
	}
}

// TransactionalConn returns an open transaction to the database.
func (d *ConnProvider) TransactionalConn(ctx context.Context) (Conn, error) {
	tx, err := d.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, realmmgr_errors.NewInternalError("failed to start transaction", err)
	}
	return &hookableConn{
		db:    tx,
		hooks: d.hooks,
	}, nil
}

// Commit handles the committing of *sql.Tx transactions, and will return an error
// if the specified connection is not transactional.
func (d *ConnProvider) Commit(conn Conn) error {
	hookableConn, ok := conn.(*hookableConn)
	if !ok {
		return realmmgr_errors.NewUnknownError("unknown Conn type", nil)
	}
	conn = hookableConn.db
	switch dbConn := conn.(type) {
	case *sql.Tx:
		if err := dbConn.Commit(); err != nil {
			return realmmgr_errors.NewInternalError("failed to commit transaction", err)
		}
		return nil
	case *sql.DB:
		return realmmgr_errors.NewInternalError("cannot commit non-transactional Conn", nil)
	default:
		return realmmgr_errors.NewUnknownError("unknown Conn type", nil)
	}
}

// Rollback rolls back an open *sql.Tx transaction, returning an error if the
// connection specified is not transactional. As this implementation uses the database/sql
// interface, it is safe to defer this rollback.
func (d *ConnProvider) Rollback(conn Conn) error {
	hookableConn, ok := conn.(*hookableConn)
	if !ok {
		return realmmgr_errors.NewUnknownError("unknown Conn type", nil)
	}
	conn = hookableConn.db
	switch dbConn := conn.(type) {
	case *sql.Tx:
		if err := dbConn.Rollback(); err != nil {
			// TODO: error wrapping
			return err
		}
		return nil
	case *sql.DB:
		return realmmgr_errors.NewInternalError("cannot rollback non-transactional Conn", nil)
	default:
		return realmmgr_errors.NewUnknownError("unknown Conn type", nil)
	}
}
