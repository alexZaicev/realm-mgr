package utils

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/alexZaicev/realm-mgr/internal/adapters/postgres/models"
	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	"github.com/alexZaicev/realm-mgr/internal/drivers/config"
	"github.com/alexZaicev/realm-mgr/internal/drivers/pgdb"
	"github.com/google/uuid"
)

const (
	ConfigDBHost    = "database.host"
	ConfigDBPort    = "database.port"
	ConfigDBUser    = "database.user"
	ConfigDBPass    = "database.password"
	ConfigDBName    = "database.name"
	ConfigDBSSLMode = "database.ssl_mode"
)

var Tables = []string{
	models.RealmTableName,
}

type DB struct {
	connProvider *pgdb.ConnProvider
}

func NewDB(cfg config.Config) (*DB, error) {
	dbHost, err := config.Get[string](cfg, ConfigDBHost)
	if err != nil {
		return nil, err
	}
	dbPort, err := config.Get[int](cfg, ConfigDBPort)
	if err != nil {
		return nil, err
	}
	dbUser, err := config.Get[string](cfg, ConfigDBUser)
	if err != nil {
		return nil, err
	}
	dbPassword, err := config.Get[string](cfg, ConfigDBPass)
	if err != nil {
		return nil, err
	}
	dbName, err := config.Get[string](cfg, ConfigDBName)
	if err != nil {
		return nil, err
	}
	dbSSLMode, err := config.Get[string](cfg, ConfigDBSSLMode)
	if err != nil {
		return nil, err
	}

	connProvider, err := pgdb.NewConnProvider(dbHost, fmt.Sprintf("%d", dbPort), dbUser, dbPassword, dbName, dbSSLMode)
	if err != nil {
		return nil, err
	}

	return &DB{
		connProvider: connProvider,
	}, nil
}

func (d *DB) Wipe() error {
	conn := d.connProvider.Conn()

	for _, table := range Tables {
		query := sq.StatementBuilder.
			PlaceholderFormat(sq.Dollar).
			Delete(table)

		if _, err := query.RunWith(conn).ExecContext(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) ExecuteInsertQueries(ctx context.Context, queries ...sq.InsertBuilder) ([]sql.Result, error) {
	conn, err := d.connProvider.TransactionalConn(ctx)
	if err != nil {
		return nil, err
	}
	defer d.connProvider.Rollback(conn)

	results := make([]sql.Result, len(queries))
	for i, query := range queries {
		results[i], err = query.RunWith(conn).ExecContext(ctx)
		if err != nil {
			return nil, err
		}
	}

	if commitErr := d.connProvider.Commit(conn); commitErr != nil {
		return nil, err
	}
	return results, nil
}

func (d *DB) RunQuery(ctx context.Context, query sq.SelectBuilder) (*sql.Rows, error) {
	conn := d.connProvider.Conn()
	rows, err := query.RunWith(conn).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DB) GetRealms(query sq.SelectBuilder) ([]*entities.Realm, error) {
	realms := make([]*entities.Realm, 0)

	rows, err := d.RunQuery(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var realm entities.Realm

		var realmKey uuid.UUID
		var dbStatus string
		var deletedAt sql.NullTime

		if scanErr := rows.Scan(
			&realmKey,
			&realm.ID,
			&realm.Name,
			&realm.Description,
			&dbStatus,
			&realm.CreatedAt,
			&realm.UpdatedAt,
			&deletedAt,
		); scanErr != nil {
			return nil, scanErr
		}

		status, ok := models.StatusDBValues[dbStatus]
		if !ok {
			return nil, fmt.Errorf("unexpected status type: %s", dbStatus)
		}
		realm.Status = status

		if deletedAt.Valid {
			realm.DeletedAt = deletedAt.Time
		}

		realms = append(realms, &realm)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return realms, nil
}
