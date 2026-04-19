package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/port"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"
)

type AdminLoginUC struct {
	admin    port.AdminRepository
	password port.PasswordHasher
	token    port.TokenGenerator
}

func NewAdminLoginUC(
	admin port.AdminRepository,
	password port.PasswordHasher,
	token port.TokenGenerator,
) *AdminLoginUC {
	return &AdminLoginUC{
		admin:    admin,
		password: password,
		token:    token,
	}
}

func (uc *AdminLoginUC) Execute(ctx context.Context, in dto.AdminLoginInput) (dto.AdminLoginOutput, error) {
	// Find the admin
	admin, err := uc.admin.GetByLogin(ctx, in.Login)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.AdminLoginOutput{}, ucerrs.ErrInvalidCredentials
		}
		return dto.AdminLoginOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetAdminByLoginDB, err,
		)
	}

	// Validate password
	if !uc.password.Compare(admin.PasswordHash(), in.Password) {
		return dto.AdminLoginOutput{}, ucerrs.ErrInvalidCredentials
	}

	// Generate the JWT token
	token, err := uc.token.Generate(admin.ID())
	if err != nil {
		return dto.AdminLoginOutput{}, ucerrs.Wrap(
			ucerrs.ErrGenerateToken, err,
		)
	}

	return dto.AdminLoginOutput{Token: token}, nil
}
