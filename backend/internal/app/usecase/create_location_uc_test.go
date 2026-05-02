package usecase_test

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/usecase"
	"backend/internal/domain/model"
	"backend/internal/domain/port/mocks"
	pkgerrs "backend/pkg/errs"
	"backend/pkg/utils"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLocationUC_Execute(t *testing.T) {
	type adapter struct {
		location     *mocks.MockLocationRepository
		item         *mocks.MockItemRepository
		locationItem *mocks.MockLocationItemRepository
	}

	type testCase struct {
		name          string
		input         dto.CreateLocationInput
		mockBehaviour func(a adapter)
		expectErr     error
	}

	testItem, _ := model.NewItem(
		"super-item",
		utils.VPtr("description for super-item"),
		"drinks",
		"https://new-photo.jpg",
		model.RestoreNutrition(
			utils.VPtr(100),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(10)),
		),
	)

	var tests = []testCase{
		{
			name: "Success",
			input: dto.CreateLocationInput{
				Slug:    "test-slug",
				Name:    "Test LocationResponse for mall",
				Address: "Brooklyn, st. main Avenue, 2378",
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Location")).Return(nil)
				a.item.EXPECT().ListAll(mock.Anything).Return([]*model.Item{testItem}, nil)
				a.locationItem.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.LocationItem")).Return(nil)
			},
			expectErr: nil,
		},
		{
			name: "Failure - invalid input",
			input: dto.CreateLocationInput{
				Slug:    "",
				Name:    "",
				Address: "",
			},
			mockBehaviour: func(a adapter) {},
			expectErr:     ucerrs.ErrInvalidInput,
		},
		{
			name: "Failure - location already exists",
			input: dto.CreateLocationInput{
				Slug:    "test-slug",
				Name:    "Test LocationResponse for mall",
				Address: "Brooklyn, st. main Avenue, 2378",
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Location")).Return(pkgerrs.ErrObjectAlreadyExists)
			},
			expectErr: ucerrs.ErrLocationAlreadyExists,
		},
		{
			name: "Failure - create location db error",
			input: dto.CreateLocationInput{
				Slug:    "test-slug",
				Name:    "Test LocationResponse for mall",
				Address: "Brooklyn, st. main Avenue, 2378",
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Location")).Return(errors.New("db error"))
			},
			expectErr: ucerrs.ErrCreateLocationDB,
		},
		{
			name: "Failure - list items db error",
			input: dto.CreateLocationInput{
				Slug:    "test-slug",
				Name:    "Test LocationResponse for mall",
				Address: "Brooklyn, st. main Avenue, 2378",
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().Create(mock.Anything, mock.AnythingOfType("*model.Location")).Return(nil)
				a.item.EXPECT().ListAll(mock.Anything).Return(nil, errors.New("db error"))
			},
			expectErr: ucerrs.ErrListAllItemsDB,
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

			uc := usecase.NewCreateLocationUC(mocks.FakeTxManager{}, locRepo, itemRepo, locItemRepo)

			out, err := uc.Execute(context.Background(), tt.input)

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, out.Location)
				assert.Equal(t, tt.input.Name, out.Location.Name)
			}
		})
	}
}
