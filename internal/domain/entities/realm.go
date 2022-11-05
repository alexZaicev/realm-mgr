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

func (r Realm) DeepCopyRealm() Realm {
	return Realm{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}
