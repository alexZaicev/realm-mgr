package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/alexZaicev/realm-mgr/internal/adapters/postgres/models"
	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
)

var selectRealmColumns = []string{
	models.RealmColumnID.WithTable(),
	models.RealmColumnName.WithTable(),
	models.RealmColumnDesc.WithTable(),
	models.RealmColumnStatus.WithTable(),
	models.RealmColumnCreatedAt.WithTable(),
	models.RealmColumnUpdatedAt.WithTable(),
	models.RealmColumnDeletedAt.WithTable(),
}

func (d *DataStore) GetRealm(ctx context.Context, realmID uuid.UUID, status entities.Status) (entities.Realm, error) {
	dbStatus, ok := models.StatusEnumValues[status]
	if !ok {
		return entities.Realm{}, realmmgr_errors.NewUnknownError(
			fmt.Sprintf("unexpected status type: %d", status),
			nil,
		)
	}

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(selectRealmColumns...).
		From(models.RealmTableName).
		Where(sq.Eq{
			models.RealmColumnID.WithTable():     realmID,
			models.RealmColumnStatus.WithTable(): dbStatus,
		})

	var realm entities.Realm

	var statusDBVal string
	var deletedAt sql.NullTime

	if err := query.RunWith(d.db).QueryRowContext(ctx).Scan(
		&realm.ID,
		&realm.Name,
		&realm.Description,
		&statusDBVal,
		&realm.CreatedAt,
		&realm.UpdatedAt,
		&deletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Realm{}, realmmgr_errors.NewNotFoundError("realm not found", err)
		}
		return entities.Realm{}, realmmgr_errors.NewInternalError("realm select failed", err)
	}

	realmStatus, ok := models.StatusDBValues[statusDBVal]
	if !ok {
		return entities.Realm{}, realmmgr_errors.NewUnknownError(
			fmt.Sprintf("unexpected status type: %s", statusDBVal),
			nil,
		)
	}
	realm.Status = realmStatus

	if deletedAt.Valid {
		realm.DeletedAt = deletedAt.Time
	}

	return realm, nil
}
