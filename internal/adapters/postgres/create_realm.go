package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/alexZaicev/realm-mgr/internal/adapters/postgres/models"
	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
)

var insertRealmColumns = []string{
	models.RealmColumnKey.String(),
	models.RealmColumnID.String(),
	models.RealmColumnName.String(),
	models.RealmColumnDesc.String(),
	models.RealmColumnStatus.String(),
	models.RealmColumnCreatedAt.String(),
	models.RealmColumnUpdatedAt.String(),
}

func (d *DataStore) CreateRealm(ctx context.Context, realm entities.Realm) error {
	key, err := d.uuidgen.New()
	if err != nil {
		return realmmgr_errors.NewInternalError("failed to generate UUID key", err)
	}

	status, ok := models.StatusEnumValues[realm.Status]
	if !ok {
		return realmmgr_errors.NewUnknownError(
			fmt.Sprintf("unexpected status type: %d", realm.Status),
			nil,
		)
	}

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert(models.RealmTableName).
		Columns(insertRealmColumns...).
		Values(
			key,
			realm.ID,
			realm.Name,
			realm.Description,
			status,
			realm.CreatedAt,
			realm.UpdatedAt,
		)

	if _, insertErr := query.RunWith(d.db).ExecContext(ctx); insertErr != nil {
		return realmmgr_errors.NewInternalError("realm insert failed", insertErr)
	}

	return nil
}
