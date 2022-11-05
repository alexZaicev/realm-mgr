package utils

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/alexZaicev/realm-mgr/internal/adapters/postgres/models"
	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	"github.com/google/uuid"
)

func GenerateRealmInsertQueries(realms ...entities.Realm) ([]sq.InsertBuilder, error) {
	queries := make([]sq.InsertBuilder, 0, len(realms))

	for _, realm := range realms {
		dbStatus, ok := models.StatusEnumValues[realm.Status]
		if !ok {
			return nil, fmt.Errorf("unexpected status type: %d", realm.Status)
		}

		query := sq.StatementBuilder.
			PlaceholderFormat(sq.Dollar).
			Insert(models.RealmTableName).
			Columns(
				models.RealmColumnKey.String(),
				models.RealmColumnID.String(),
				models.RealmColumnName.String(),
				models.RealmColumnDesc.String(),
				models.RealmColumnStatus.String(),
				models.RealmColumnCreatedAt.String(),
				models.RealmColumnUpdatedAt.String(),
			).
			Values(
				uuid.New(),
				realm.ID,
				realm.Name,
				realm.Description,
				dbStatus,
				realm.CreatedAt,
				realm.UpdatedAt,
			)
		queries = append(queries, query)
	}

	return queries, nil
}

func GetRealmsQuery() sq.SelectBuilder {
	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(
			models.RealmColumnKey.String(),
			models.RealmColumnID.String(),
			models.RealmColumnName.String(),
			models.RealmColumnDesc.String(),
			models.RealmColumnStatus.String(),
			models.RealmColumnCreatedAt.String(),
			models.RealmColumnUpdatedAt.String(),
			models.RealmColumnDeletedAt.String(),
		).
		From(models.RealmTableName)

	return query
}
