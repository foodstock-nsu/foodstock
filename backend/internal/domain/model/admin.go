package model

import (
	pkgerrs "backend/pkg/errs"
	"time"

	"github.com/google/uuid"
)

// ================ Rich model for Admin ================

type Admin struct {
	id           uuid.UUID
	login        string
	passwordHash string
	createdAt    time.Time
}

func NewAdmin(login, passwordHash string) (*Admin, error) {
	if len(login) == 0 {
		return nil, pkgerrs.NewValueRequiredError("login")
	}
	if len(login) <= 4 {
		return nil, pkgerrs.NewValueInvalidError("login")
	}

	if len(passwordHash) == 0 {
		return nil, pkgerrs.NewValueRequiredError("password_hash")
	}

	return &Admin{
		id:           uuid.New(),
		login:        login,
		passwordHash: passwordHash,
		createdAt:    time.Now().UTC(),
	}, nil
}

func RestoreAdmin(
	id uuid.UUID,
	login, passwordHash string,
	createdAt time.Time,
) *Admin {
	return &Admin{
		id:           id,
		login:        login,
		passwordHash: passwordHash,
		createdAt:    createdAt,
	}
}

// ================ Read-Only ================

func (a *Admin) ID() uuid.UUID        { return a.id }
func (a *Admin) Login() string        { return a.login }
func (a *Admin) PasswordHash() string { return a.passwordHash }
func (a *Admin) CreatedAt() time.Time { return a.createdAt }
