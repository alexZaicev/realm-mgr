package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/alexZaicev/realm-mgr/internal/adapters/postgres/models"
	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
)

func (d *DataStore) DeleteRealm(ctx context.Context, realmID uuid.UUID, statuses ...entities.Status) error {
	var dbStatuses []string
	for _, status := range statuses {
		dbStatus, ok := models.StatusEnumValues[status]
		if !ok {
			return realmmgr_errors.NewUnknownError(
				fmt.Sprintf("unexpected status type: %d", status),
				nil,
			)
		}
		dbStatuses = append(dbStatuses, dbStatus)
	}

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Delete(models.RealmTableName).
		Where(sq.Eq{
			models.RealmColumnID.String():     realmID,
			models.RealmColumnStatus.String(): dbStatuses,
		})

	if _, err := query.RunWith(d.db).ExecContext(ctx); err != nil {
		return realmmgr_errors.NewInternalError("realm delete failed", err)
	}

	return nil
}
