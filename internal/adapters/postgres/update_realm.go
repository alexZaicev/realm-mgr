package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/alexZaicev/realm-mgr/internal/adapters/postgres/models"
	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
)

func (d *DataStore) UpdateRealm(ctx context.Context, realm entities.Realm, currentStatus entities.Status) error {
	dbStatus, ok := models.StatusEnumValues[currentStatus]
	if !ok {
		return realmmgr_errors.NewUnknownError(
			fmt.Sprintf("unexpected status type: %d", currentStatus),
			nil,
		)
	}

	realmStatus, ok := models.StatusEnumValues[realm.Status]
	if !ok {
		return realmmgr_errors.NewUnknownError(
			fmt.Sprintf("unexpected status type: %d", realm.Status),
			nil,
		)
	}

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update(models.RealmTableName).
		SetMap(map[string]interface{}{
			models.RealmColumnName.String():      realm.Name,
			models.RealmColumnDesc.String():      realm.Description,
			models.RealmColumnUpdatedAt.String(): realm.UpdatedAt,
			models.RealmColumnStatus.String():    realmStatus,
		}).
		Where(sq.Eq{
			models.RealmColumnID.String():     realm.ID,
			models.RealmColumnStatus.String(): dbStatus,
		})

	if _, err := query.RunWith(d.db).ExecContext(ctx); err != nil {
		return realmmgr_errors.NewInternalError("realm update failed", err)
	}

	return nil
}
