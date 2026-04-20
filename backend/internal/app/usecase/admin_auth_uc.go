package usecase

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/domain/port"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"
)

type AdminAuthUC struct {
	admin    port.AdminRepository
	password port.PasswordHasher
	token    port.TokenGenerator
}

func NewAdminAuthUC(
	admin port.AdminRepository,
	password port.PasswordHasher,
	token port.TokenGenerator,
) *AdminAuthUC {
	return &AdminAuthUC{
		admin:    admin,
		password: password,
		token:    token,
	}
}

func (uc *AdminAuthUC) Execute(ctx context.Context, in dto.AdminAuthInput) (dto.AdminAuthOutput, error) {
	// Find the admin
	admin, err := uc.admin.GetByLogin(ctx, in.Login)
	if err != nil {
		if errors.Is(err, pkgerrs.ErrObjectNotFound) {
			return dto.AdminAuthOutput{}, ucerrs.ErrInvalidCredentials
		}
		return dto.AdminAuthOutput{}, ucerrs.Wrap(
			ucerrs.ErrGetAdminByLoginDB, err,
		)
	}

	// Validate password
	if !uc.password.Compare(admin.PasswordHash(), in.Password) {
		return dto.AdminAuthOutput{}, ucerrs.ErrInvalidCredentials
	}

	// Generate the JWT token
	token, err := uc.token.Generate(admin.ID())
	if err != nil {
		return dto.AdminAuthOutput{}, ucerrs.Wrap(
			ucerrs.ErrGenerateToken, err,
		)
	}

	return dto.AdminAuthOutput{Token: token}, nil
}
