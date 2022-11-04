package uuidgenerator

import "github.com/google/uuid"

// Generator defines an interface for creating new UUIDs.
type Generator interface {
	New() (uuid.UUID, error)
}

// GoogleUUIDGenerator provides a real implementation of the Generator interface
// to be used in production using Google's uuid package.
type GoogleUUIDGenerator struct{}

// NewGoogleUUIDGenerator returns a GoogleUUIDGenerator, and can be used in dependency
// injection frameworks.
func NewGoogleUUIDGenerator() GoogleUUIDGenerator {
	return GoogleUUIDGenerator{}
}

// New generates a UUID.
func (g GoogleUUIDGenerator) New() (uuid.UUID, error) {
	return uuid.NewRandom()
}
