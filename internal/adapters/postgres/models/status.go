package models

import "github.com/alexZaicev/realm-mgr/internal/domain/entities"

var (
	StatusEnumValues = map[entities.Status]string{
		entities.StatusActive:   "active",
		entities.StatusDraft:    "draft",
		entities.StatusDisabled: "disabled",
		entities.StatusDeleted:  "deleted",
	}

	StatusDBValues = func() map[string]entities.Status {
		result := make(map[string]entities.Status)
		for k, v := range StatusEnumValues {
			result[v] = k
		}
		return result
	}()
)
