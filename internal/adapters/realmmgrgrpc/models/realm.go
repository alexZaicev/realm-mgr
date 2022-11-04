package models

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
)

func RealmFromDomain(realm entities.Realm) (*realm_mgr_v1.Realm, error) {
	status, ok := StatusEnumValues[realm.Status]
	if !ok {
		return nil, realmmgr_errors.NewUnknownError(fmt.Sprintf("unexpected status type: %d", realm.Status), nil)
	}

	return &realm_mgr_v1.Realm{
		Id:          realm.ID.String(),
		Name:        realm.Name,
		Description: realm.Description,
		Status:      status,
		CreatedAt:   timestamppb.New(realm.CreatedAt),
		UpdatedAt:   timestamppb.New(realm.UpdatedAt),
	}, nil
}
