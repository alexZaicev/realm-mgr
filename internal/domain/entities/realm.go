package entities

import (
	"time"

	"github.com/google/uuid"
)

type Realm struct {
	ID          uuid.UUID
	Name        string
	Description string
	Status      Status

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (r Realm) Merge(realm Realm) Realm {
	r.Name = realm.Name
	r.Description = realm.Description
	r.UpdatedAt = realm.UpdatedAt

	return r
}
