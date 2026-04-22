package usecase_test

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/usecase"
	"backend/internal/domain/port/mocks"
	pkgerrs "backend/pkg/errs"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteLocationUC_Execute(t *testing.T) {
	type adapter struct {
		location     *mocks.MockLocationRepository
		locationItem *mocks.MockLocationItemRepository
	}

	type testCase struct {
		name          string
		input         dto.DeleteLocationInput
		mockBehaviour func(a adapter)
		expectErr     error
	}

	testID := uuid.New()

	var tests = []testCase{
		{
			name: "Success",
			input: dto.DeleteLocationInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().Delete(mock.Anything, testID).Return(nil)
				a.locationItem.EXPECT().DeleteByLocationID(mock.Anything, testID).Return(nil)
			},
			expectErr: nil,
		},
		{
			name: "Failure - location not found",
			input: dto.DeleteLocationInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().Delete(mock.Anything, testID).Return(pkgerrs.ErrObjectNotFound)
			},
			expectErr: ucerrs.ErrLocationNotFound,
		},
		{
			name: "Failure - delete location db error",
			input: dto.DeleteLocationInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().Delete(mock.Anything, testID).Return(errors.New("db error"))
			},
			expectErr: ucerrs.ErrDeleteLocationDB,
		},
		{
			name: "Failure - delete location items db error",
			input: dto.DeleteLocationInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().Delete(mock.Anything, testID).Return(nil)
				a.locationItem.EXPECT().DeleteByLocationID(mock.Anything, testID).Return(errors.New("db error"))
			},
			expectErr: ucerrs.ErrDeleteLocationItemByLocationIDDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locRepo := mocks.NewMockLocationRepository(t)
			locItemRepo := mocks.NewMockLocationItemRepository(t)

			tt.mockBehaviour(adapter{
				location:     locRepo,
				locationItem: locItemRepo,
			})

			uc := usecase.NewDeleteLocationUC(mocks.FakeTxManager{}, locRepo, locItemRepo)

			err := uc.Execute(context.Background(), tt.input)

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
