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

func TestDeleteItemUC_Execute(t *testing.T) {
	type adapter struct {
		item         *mocks.MockItemRepository
		locationItem *mocks.MockLocationItemRepository
	}

	type testCase struct {
		name          string
		input         dto.DeleteItemInput
		mockBehaviour func(a adapter)
		expectErr     error
	}

	testID := uuid.New()

	var tests = []testCase{
		{
			name: "Success",
			input: dto.DeleteItemInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.item.EXPECT().Delete(mock.Anything, testID).Return(nil)
				a.locationItem.EXPECT().DeleteByItemID(mock.Anything, testID).Return(nil)
			},
			expectErr: nil,
		},
		{
			name: "Failure - item not found",
			input: dto.DeleteItemInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.item.EXPECT().Delete(mock.Anything, testID).Return(pkgerrs.ErrObjectNotFound)
			},
			expectErr: ucerrs.ErrItemNotFound,
		},
		{
			name: "Failure - delete item db error",
			input: dto.DeleteItemInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.item.EXPECT().Delete(mock.Anything, testID).Return(errors.New("db error"))
			},
			expectErr: ucerrs.ErrDeleteItemDB,
		},
		{
			name: "Failure - delete location items db error",
			input: dto.DeleteItemInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.item.EXPECT().Delete(mock.Anything, testID).Return(nil)
				a.locationItem.EXPECT().DeleteByItemID(mock.Anything, testID).Return(errors.New("db error"))
			},
			expectErr: ucerrs.ErrDeleteLocationItemsByItemIDDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			itemRepo := mocks.NewMockItemRepository(t)
			locItemRepo := mocks.NewMockLocationItemRepository(t)

			tt.mockBehaviour(adapter{
				item:         itemRepo,
				locationItem: locItemRepo,
			})

			uc := usecase.NewDeleteItemUC(mocks.FakeTxManager{}, itemRepo, locItemRepo)

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
