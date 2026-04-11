package model

import (
	"errors"
	"time"

	pkgerrs "backend/pkg/errs"

	"github.com/google/uuid"
)

var (
	ErrCannotActivate   = errors.New("location is already activated")
	ErrCannotDeactivate = errors.New("location is already deactivated")
)

// ================ Rich model for Location (e.g. Fridge) ================

type Location struct {
	id        uuid.UUID
	slug      string // mark from QR-code, for example "nsu_1"
	name      string
	address   string
	isActive  bool
	createdAt time.Time
}

func NewLocation(slug, name, address string) (*Location, error) {
	if len(slug) < 4 {
		return nil, pkgerrs.NewValueInvalidError("slug")
	}
	if len(name) < 4 {
		return nil, pkgerrs.NewValueInvalidError("name")
	}
	if len(address) < 20 {
		return nil, pkgerrs.NewValueInvalidError("address")
	}

	return &Location{
		id:        uuid.New(),
		slug:      slug,
		name:      name,
		address:   address,
		isActive:  true,
		createdAt: time.Now().UTC(),
	}, nil
}

func RestoreLocation(
	id uuid.UUID,
	slug, name, address string,
	isActive bool,
	createdAt time.Time,
) *Location {
	return &Location{
		id:        id,
		slug:      slug,
		name:      name,
		address:   address,
		isActive:  isActive,
		createdAt: createdAt,
	}
}

// ================ Read-Only ================

func (l *Location) ID() uuid.UUID        { return l.id }
func (l *Location) Slug() string         { return l.slug }
func (l *Location) Name() string         { return l.name }
func (l *Location) Address() string      { return l.address }
func (l *Location) IsActive() bool       { return l.isActive }
func (l *Location) CreatedAt() time.Time { return l.createdAt }

// ================ Business Logic ================

func (l *Location) IsOperational() bool {
	return l.isActive
}

func (l *Location) ValidateQRCode(code string) bool {
	return l.slug == code && l.isActive
}

// ================ Mutation ================

func (l *Location) Update(slug, name, address *string) error {
	if slug != nil && len(*slug) < 4 {
		return pkgerrs.NewValueInvalidError("slug")
	}
	if name != nil && len(*name) < 4 {
		return pkgerrs.NewValueInvalidError("name")
	}
	if address != nil && len(*address) < 20 {
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

func (l *Location) Activate() error {
	if l.isActive {
		return ErrCannotActivate
	}
	l.isActive = true
	return nil
}

func (l *Location) Deactivate() error {
	if !l.isActive {
		return ErrCannotDeactivate
	}
	l.isActive = false
	return nil
}
