package models

import (
	"github.com/alexZaicev/realm-mgr/internal/domain/entities"
	realm_mgr_v1 "github.com/alexZaicev/realm-mgr/proto/go/realm_mgr/v1"
)

const (
	InternalErrMsg = "an internal error occurred"
	InvalidRealmID = "realm ID was not a valid UUID: %s"
)

var (
	StatusEnumValues = map[entities.Status]realm_mgr_v1.EnumStatus{
		entities.StatusActive:   realm_mgr_v1.EnumStatus_ENUM_STATUS_ACTIVE,
		entities.StatusDraft:    realm_mgr_v1.EnumStatus_ENUM_STATUS_DRAFT,
		entities.StatusDisabled: realm_mgr_v1.EnumStatus_ENUM_STATUS_DISABLED,
		entities.StatusDeleted:  realm_mgr_v1.EnumStatus_ENUM_STATUS_DELETED,
	}

	StatusGRPCValues = func() map[realm_mgr_v1.EnumStatus]entities.Status {
		result := make(map[realm_mgr_v1.EnumStatus]entities.Status)
		for k, v := range StatusEnumValues {
			result[v] = k
		}
		return result
	}()
)
