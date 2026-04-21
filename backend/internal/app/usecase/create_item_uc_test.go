package usecase_test

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/usecase"
	"backend/internal/domain/model"
	"backend/internal/domain/port/mocks"
	"backend/pkg/utils"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateItemUC_Execute(t *testing.T) {
	type adapter struct {
		location     *mocks.MockLocationRepository
		item         *mocks.MockItemRepository
		locationItem *mocks.MockLocationItemRepository
	}

	type testCase struct {
		name          string
		input         dto.CreateItemInput
		mockBehaviour func(a adapter)
		expectErr     error
	}

	validInput := dto.CreateItemInput{
		Name:        "Valid ItemOutput Name",
		Description: utils.VPtr("Valid Description"),
		Category:    "drinks",
		PhotoURL:    "https://example.com/photo.jpg",
		Nutrition: &dto.NutritionOutput{
			Calories: utils.VPtr(250),
			Proteins: utils.VPtr(float64(10.5)),
			Fats:     utils.VPtr(float64(5.0)),
			Carbs:    utils.VPtr(float64(30.0)),
		},
	}

	testLocation, _ := model.NewLocation(
		"test-slug",
		"Test LocationOutput for mall",
		"Brooklyn, st. main Avenue, 2378",
	)

	var tests = []testCase{
		{
			name:  "Success",
			input: validInput,
			mockBehaviour: func(a adapter) {
				a.item.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.ItemOutput")).Return(nil)
				a.location.EXPECT().List(mock.Anything).Return([]*model.Location{testLocation}, nil)
				a.locationItem.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.LocationItem")).Return(nil)
			},
			expectErr: nil,
		},
		{
			name: "Failure - invalid input",
			input: dto.CreateItemInput{
				Name: "Bad",
			},
			mockBehaviour: func(a adapter) {},
			expectErr:     ucerrs.ErrInvalidInput,
		},
		{
			name:  "Failure - item create db error",
			input: validInput,
			mockBehaviour: func(a adapter) {
				a.item.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.ItemOutput")).Return(errors.New("db error"))
			},
			expectErr: ucerrs.ErrCreateItemDB,
		},
		{
			name:  "Failure - list locations db error",
			input: validInput,
			mockBehaviour: func(a adapter) {
				a.item.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.ItemOutput")).Return(nil)
				a.location.EXPECT().List(mock.Anything).Return(nil, errors.New("db error"))
			},
			expectErr: ucerrs.ErrListLocationsDB,
		},
		{
			name:  "Failure - create location item db error",
			input: validInput,
			mockBehaviour: func(a adapter) {
				a.item.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.ItemOutput")).Return(nil)
				a.location.EXPECT().List(mock.Anything).Return([]*model.Location{testLocation}, nil)
				a.locationItem.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.LocationItem")).Return(errors.New("db error"))
			},
			expectErr: ucerrs.ErrCreateLocationItemDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locRepo := mocks.NewMockLocationRepository(t)
			itemRepo := mocks.NewMockItemRepository(t)
			locItemRepo := mocks.NewMockLocationItemRepository(t)

			tt.mockBehaviour(adapter{
				location:     locRepo,
				item:         itemRepo,
				locationItem: locItemRepo,
			})

			uc := usecase.NewCreateItemUC(mocks.FakeTxManager{}, locRepo, itemRepo, locItemRepo)

			out, err := uc.Execute(context.Background(), tt.input)

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, out.Item)
				assert.Equal(t, tt.input.Name, out.Item.Name)
			}
		})
	}
}
