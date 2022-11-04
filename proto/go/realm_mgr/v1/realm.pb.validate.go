// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: realm_mgr/v1/realm.proto

package realm_mgr_v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// define the regex for a UUID once up-front
var _realm_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on Realm with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Realm) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Realm with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in RealmMultiError, or nil if none found.
func (m *Realm) ValidateAll() error {
	return m.validate(true)
}

func (m *Realm) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetId()); err != nil {
		err = RealmValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := RealmValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Description

	// no validation rules for Status

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RealmValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RealmValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RealmValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RealmValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RealmValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RealmValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return RealmMultiError(errors)
	}

	return nil
}

func (m *Realm) _validateUuid(uuid string) error {
	if matched := _realm_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// RealmMultiError is an error wrapping multiple validation errors returned by
// Realm.ValidateAll() if the designated constraints aren't met.
type RealmMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RealmMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RealmMultiError) AllErrors() []error { return m }

// RealmValidationError is the validation error returned by Realm.Validate if
// the designated constraints aren't met.
type RealmValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RealmValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RealmValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RealmValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RealmValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RealmValidationError) ErrorName() string { return "RealmValidationError" }

// Error satisfies the builtin error interface
func (e RealmValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRealm.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RealmValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RealmValidationError{}

// Validate checks the field values on GetRealmRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetRealmRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetRealmRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetRealmRequestMultiError, or nil if none found.
func (m *GetRealmRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetRealmRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetId()); err != nil {
		err = GetRealmRequestValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Status

	if len(errors) > 0 {
		return GetRealmRequestMultiError(errors)
	}

	return nil
}

func (m *GetRealmRequest) _validateUuid(uuid string) error {
	if matched := _realm_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// GetRealmRequestMultiError is an error wrapping multiple validation errors
// returned by GetRealmRequest.ValidateAll() if the designated constraints
// aren't met.
type GetRealmRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetRealmRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetRealmRequestMultiError) AllErrors() []error { return m }

// GetRealmRequestValidationError is the validation error returned by
// GetRealmRequest.Validate if the designated constraints aren't met.
type GetRealmRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetRealmRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetRealmRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetRealmRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetRealmRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetRealmRequestValidationError) ErrorName() string { return "GetRealmRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetRealmRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetRealmRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetRealmRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetRealmRequestValidationError{}

// Validate checks the field values on GetRealmResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetRealmResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetRealmResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetRealmResponseMultiError, or nil if none found.
func (m *GetRealmResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetRealmResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetRealm()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetRealmResponseValidationError{
					field:  "Realm",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetRealmResponseValidationError{
					field:  "Realm",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRealm()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetRealmResponseValidationError{
				field:  "Realm",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetRealmResponseMultiError(errors)
	}

	return nil
}

// GetRealmResponseMultiError is an error wrapping multiple validation errors
// returned by GetRealmResponse.ValidateAll() if the designated constraints
// aren't met.
type GetRealmResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetRealmResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetRealmResponseMultiError) AllErrors() []error { return m }

// GetRealmResponseValidationError is the validation error returned by
// GetRealmResponse.Validate if the designated constraints aren't met.
type GetRealmResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetRealmResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetRealmResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetRealmResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetRealmResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetRealmResponseValidationError) ErrorName() string { return "GetRealmResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetRealmResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetRealmResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetRealmResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetRealmResponseValidationError{}

// Validate checks the field values on CreateRealmRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateRealmRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateRealmRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateRealmRequestMultiError, or nil if none found.
func (m *CreateRealmRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateRealmRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := CreateRealmRequestValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Description

	if len(errors) > 0 {
		return CreateRealmRequestMultiError(errors)
	}

	return nil
}

// CreateRealmRequestMultiError is an error wrapping multiple validation errors
// returned by CreateRealmRequest.ValidateAll() if the designated constraints
// aren't met.
type CreateRealmRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateRealmRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateRealmRequestMultiError) AllErrors() []error { return m }

// CreateRealmRequestValidationError is the validation error returned by
// CreateRealmRequest.Validate if the designated constraints aren't met.
type CreateRealmRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateRealmRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateRealmRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateRealmRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateRealmRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateRealmRequestValidationError) ErrorName() string {
	return "CreateRealmRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateRealmRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateRealmRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateRealmRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateRealmRequestValidationError{}

// Validate checks the field values on CreateRealmResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateRealmResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateRealmResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateRealmResponseMultiError, or nil if none found.
func (m *CreateRealmResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateRealmResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetRealm()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateRealmResponseValidationError{
					field:  "Realm",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateRealmResponseValidationError{
					field:  "Realm",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRealm()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateRealmResponseValidationError{
				field:  "Realm",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CreateRealmResponseMultiError(errors)
	}

	return nil
}

// CreateRealmResponseMultiError is an error wrapping multiple validation
// errors returned by CreateRealmResponse.ValidateAll() if the designated
// constraints aren't met.
type CreateRealmResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateRealmResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateRealmResponseMultiError) AllErrors() []error { return m }

// CreateRealmResponseValidationError is the validation error returned by
// CreateRealmResponse.Validate if the designated constraints aren't met.
type CreateRealmResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateRealmResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateRealmResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateRealmResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateRealmResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateRealmResponseValidationError) ErrorName() string {
	return "CreateRealmResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateRealmResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateRealmResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateRealmResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateRealmResponseValidationError{}

// Validate checks the field values on ReleaseRealmRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ReleaseRealmRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReleaseRealmRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ReleaseRealmRequestMultiError, or nil if none found.
func (m *ReleaseRealmRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ReleaseRealmRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetId()); err != nil {
		err = ReleaseRealmRequestValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ReleaseRealmRequestMultiError(errors)
	}

	return nil
}

func (m *ReleaseRealmRequest) _validateUuid(uuid string) error {
	if matched := _realm_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// ReleaseRealmRequestMultiError is an error wrapping multiple validation
// errors returned by ReleaseRealmRequest.ValidateAll() if the designated
// constraints aren't met.
type ReleaseRealmRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReleaseRealmRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReleaseRealmRequestMultiError) AllErrors() []error { return m }

// ReleaseRealmRequestValidationError is the validation error returned by
// ReleaseRealmRequest.Validate if the designated constraints aren't met.
type ReleaseRealmRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReleaseRealmRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReleaseRealmRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReleaseRealmRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReleaseRealmRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReleaseRealmRequestValidationError) ErrorName() string {
	return "ReleaseRealmRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ReleaseRealmRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReleaseRealmRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReleaseRealmRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReleaseRealmRequestValidationError{}

// Validate checks the field values on ReleaseRealmResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ReleaseRealmResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReleaseRealmResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ReleaseRealmResponseMultiError, or nil if none found.
func (m *ReleaseRealmResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ReleaseRealmResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetRealm()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ReleaseRealmResponseValidationError{
					field:  "Realm",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ReleaseRealmResponseValidationError{
					field:  "Realm",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRealm()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ReleaseRealmResponseValidationError{
				field:  "Realm",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ReleaseRealmResponseMultiError(errors)
	}

	return nil
}

// ReleaseRealmResponseMultiError is an error wrapping multiple validation
// errors returned by ReleaseRealmResponse.ValidateAll() if the designated
// constraints aren't met.
type ReleaseRealmResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReleaseRealmResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReleaseRealmResponseMultiError) AllErrors() []error { return m }

// ReleaseRealmResponseValidationError is the validation error returned by
// ReleaseRealmResponse.Validate if the designated constraints aren't met.
type ReleaseRealmResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReleaseRealmResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReleaseRealmResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReleaseRealmResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReleaseRealmResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReleaseRealmResponseValidationError) ErrorName() string {
	return "ReleaseRealmResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ReleaseRealmResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReleaseRealmResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReleaseRealmResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReleaseRealmResponseValidationError{}