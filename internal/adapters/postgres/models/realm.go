package models

import "fmt"

type RealmColumn string

func (c RealmColumn) String() string {
	return string(c)
}

func (c RealmColumn) WithTable() string {
	return fmt.Sprintf("%s.%s", RealmTableName, c)
}

const (
	RealmTableName = "realms"

	RealmColumnKey       RealmColumn = "key"
	RealmColumnID        RealmColumn = "id"
	RealmColumnName      RealmColumn = "name"
	RealmColumnDesc      RealmColumn = "description"
	RealmColumnStatus    RealmColumn = "status"
	RealmColumnCreatedAt RealmColumn = "created_at"
	RealmColumnUpdatedAt RealmColumn = "updated_at"
	RealmColumnDeletedAt RealmColumn = "deleted_at"
)
