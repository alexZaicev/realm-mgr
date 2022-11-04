package entities

type Status int

const (
	StatusActive = iota + 1
	StatusDraft
	StatusDisabled
	StatusDeleted
)
