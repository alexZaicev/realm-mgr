package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
)

type RealmRepository interface {
	GetRealm(ctx context.Context, realmID uuid.UUID, status entities.Status) (entities.Realm, error)
	CreateRealm(ctx context.Context, realm entities.Realm) error
	UpdateRealm(ctx context.Context, realm entities.Realm, currentStatus entities.Status) error
	DeleteRealm(ctx context.Context, realmID uuid.UUID, statuses ...entities.Status) error
}
