package models

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"

	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realmmgr_errors "github.com/alexZaicev/realm-mgr/internal/domain/errors"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
)

func RealmFromDomain(realm entities.Realm) (*realm_mgr_v1.Realm, error) {
	realmStatus, ok := StatusEnumValues[realm.Status]
	if !ok {
		return nil, realmmgr_errors.NewUnknownError(fmt.Sprintf("unexpected status type: %d", realm.Status), nil)
	}

	return &realm_mgr_v1.Realm{
		Id:          realm.ID.String(),
		Name:        realm.Name,
		Description: realm.Description,
		Status:      realmStatus,
		CreatedAt:   timestamppb.New(realm.CreatedAt),
		UpdatedAt:   timestamppb.New(realm.UpdatedAt),
	}, nil
}

func RealmToDomain(pbRealm *realm_mgr_v1.Realm) (entities.Realm, error) {
	if pbRealm == nil {
		return entities.Realm{}, realmmgr_errors.NewInvalidArgumentError("realm", realmmgr_errors.ErrMsgCannotBeNil)
	}

	realmID, err := uuid.Parse(pbRealm.Id)
	if err != nil {
		return entities.Realm{}, status.Errorf(codes.InvalidArgument, fmt.Sprintf("realm ID was not a valid UUID: %s", pbRealm.Id))
	}

	return entities.Realm{
		ID:          realmID,
		Name:        pbRealm.Name,
		Description: pbRealm.Description,
		Status:      entities.StatusDraft,
	}, nil
}
