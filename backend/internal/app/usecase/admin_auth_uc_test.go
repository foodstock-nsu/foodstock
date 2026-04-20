package usecase_test

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/usecase"
	"backend/internal/domain/model"
	"backend/internal/domain/port/mocks"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdminLoginUC_Execute(t *testing.T) {
	type adapter struct {
		admin    *mocks.MockAdminRepository
		password *mocks.MockPasswordHasher
		token    *mocks.MockTokenGenerator
	}

	type testCase struct {
		name          string
		input         dto.AdminAuthInput
		mockBehaviour func(a adapter)
		expectErr     error
	}

	testAdmin := &model.Admin{}

	var tests = []testCase{
		{
			name: "Success",
			input: dto.AdminAuthInput{
				Login:    "admin",
				Password: "password",
			},
			mockBehaviour: func(a adapter) {
				a.admin.EXPECT().GetByLogin(mock.Anything, "admin").Return(testAdmin, nil)
				a.password.EXPECT().Compare(mock.Anything, "password").Return(true)
				a.token.EXPECT().Generate(mock.Anything).Return("jwt-token", nil)
			},
			expectErr: nil,
		},
		{
			name: "Failure - admin not found",
			input: dto.AdminAuthInput{
				Login:    "unknown",
				Password: "password",
			},
			mockBehaviour: func(a adapter) {
				a.admin.EXPECT().GetByLogin(mock.Anything, "unknown").Return(nil, pkgerrs.ErrObjectNotFound)
			},
			expectErr: ucerrs.ErrInvalidCredentials,
		},
		{
			name: "Failure - db error",
			input: dto.AdminAuthInput{
				Login:    "admin",
				Password: "password",
			},
			mockBehaviour: func(a adapter) {
				a.admin.EXPECT().GetByLogin(mock.Anything, "admin").Return(nil, errors.New("db error"))
			},
			expectErr: ucerrs.ErrGetAdminByLoginDB,
		},
		{
			name: "Failure - invalid password",
			input: dto.AdminAuthInput{
				Login:    "admin",
				Password: "wrong",
			},
			mockBehaviour: func(a adapter) {
				a.admin.EXPECT().GetByLogin(mock.Anything, "admin").Return(testAdmin, nil)
				a.password.EXPECT().Compare(mock.Anything, "wrong").Return(false)
			},
			expectErr: ucerrs.ErrInvalidCredentials,
		},
		{
			name: "Failure - token generation error",
			input: dto.AdminAuthInput{
				Login:    "admin",
				Password: "password",
			},
			mockBehaviour: func(a adapter) {
				a.admin.EXPECT().GetByLogin(mock.Anything, "admin").Return(testAdmin, nil)
				a.password.EXPECT().Compare(mock.Anything, "password").Return(true)
				a.token.EXPECT().Generate(mock.Anything).Return("", errors.New("token error"))
			},
			expectErr: ucerrs.ErrGenerateToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminRepo := mocks.NewMockAdminRepository(t)
			passwordHasher := mocks.NewMockPasswordHasher(t)
			tokenGen := mocks.NewMockTokenGenerator(t)

			tt.mockBehaviour(adapter{
				admin:    adminRepo,
				password: passwordHasher,
				token:    tokenGen,
			})

			uc := usecase.NewAdminAuthUC(adminRepo, passwordHasher, tokenGen)

			out, err := uc.Execute(context.Background(), tt.input)

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "jwt-token", out.Token)
			}
		})
	}
}
