package model

import (
	"backend/pkg/utils"
	"errors"
	"time"
	"unicode/utf8"

	pkgerrs "backend/pkg/errs"

	"github.com/google/uuid"
)

var (
	ErrCannotDelete = errors.New("location is already deleted")
)

// ================ Rich model for Location (e.g. Fridge) ================

type Location struct {
	id        uuid.UUID
	slug      string // mark from QR-code, for example "nsu_1"
	name      string
	address   string
	isActive  bool
	createdAt time.Time
	deletedAt *time.Time
}

func NewLocation(slug, name, address string) (*Location, error) {
	if utf8.RuneCountInString(slug) < 4 {
		return nil, pkgerrs.NewValueInvalidError("slug")
	}
	if utf8.RuneCountInString(name) < 4 {
		return nil, pkgerrs.NewValueInvalidError("name")
	}
	if utf8.RuneCountInString(address) < 20 {
		return nil, pkgerrs.NewValueInvalidError("address")
	}

	return &Location{
		id:        uuid.New(),
		slug:      slug,
		name:      name,
		address:   address,
		isActive:  true,
		createdAt: time.Now().UTC(),
		deletedAt: nil,
	}, nil
}

func RestoreLocation(
	id uuid.UUID,
	slug, name, address string,
	isActive bool,
	createdAt time.Time,
	deletedAt *time.Time,
) *Location {
	return &Location{
		id:        id,
		slug:      slug,
		name:      name,
		address:   address,
		isActive:  isActive,
		createdAt: createdAt,
		deletedAt: deletedAt,
	}
}

// ================ Read-Only ================

func (l *Location) ID() uuid.UUID         { return l.id }
func (l *Location) Slug() string          { return l.slug }
func (l *Location) Name() string          { return l.name }
func (l *Location) Address() string       { return l.address }
func (l *Location) IsActive() bool        { return l.isActive }
func (l *Location) CreatedAt() time.Time  { return l.createdAt }
func (l *Location) DeletedAt() *time.Time { return l.deletedAt }

// ================ Business Logic ================

func (l *Location) IsOperational() bool {
	return l.isActive
}

func (l *Location) IsDeleted() bool { return l.deletedAt != nil }

func (l *Location) GetQRData(baseURL string) string {
	return baseURL + "?location_id=" + l.slug
}

func (l *Location) ValidateQRCode(slug string) bool {
	return l.slug == slug && l.isActive
}

// ================ Mutation ================

func (l *Location) Update(slug, name, address *string) error {
	if slug != nil && utf8.RuneCountInString(*slug) < 4 {
		return pkgerrs.NewValueInvalidError("slug")
	}
	if name != nil && utf8.RuneCountInString(*name) < 4 {
		return pkgerrs.NewValueInvalidError("name")
	}
	if address != nil && utf8.RuneCountInString(*address) < 20 {
		return pkgerrs.NewValueInvalidError("address")
	}

	if slug != nil {
		l.slug = *slug
	}
	if name != nil {
		l.name = *name
	}
	if address != nil {
		l.address = *address
	}

	return nil
}

func (l *Location) Activate() {
	l.isActive = true
}

func (l *Location) Deactivate() {
	l.isActive = false
}

func (l *Location) Delete() error {
	if l.DeletedAt() != nil {
		return ErrCannotDelete
	}
	l.deletedAt = utils.VPtr(time.Now().UTC())
	return nil
}
