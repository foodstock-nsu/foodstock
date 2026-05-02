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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateLocationUC_Execute(t *testing.T) {
	type adapter struct {
		location *mocks.MockLocationRepository
	}

	type testCase struct {
		name          string
		input         dto.UpdateLocationInput
		mockBehaviour func(a adapter)
		expectErr     error
	}

	testID := uuid.New()

	var tests = []testCase{
		{
			name: "Success",
			input: dto.UpdateLocationInput{
				ID:       testID,
				Slug:     utils.VPtr("updated-slug"),
				Name:     utils.VPtr("Updated Name Name Name"),
				Address:  utils.VPtr("Updated Address Address "),
				IsActive: nil,
			},
			mockBehaviour: func(a adapter) {
				testLocation, _ := model.NewLocation(
					"test-slug",
					"Test LocationDTO for mall",
					"Brooklyn, st. main Avenue, 2378",
				)
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocation, nil)
				a.location.EXPECT().Update(mock.Anything, mock.AnythingOfType("*model.Location")).Return(nil)
			},
			expectErr: nil,
		},
		//{
		//	name: "Success - activate location",
		//	input: dto.UpdateLocationInput{
		//		ID:       testID,
		//		Slug:     utils.VPtr("updated-slug"),
		//		Name:     utils.VPtr("Updated Name"),
		//		Address:  utils.VPtr("Updated Address address"),
		//		IsActive: utils.VPtr(true),
		//	},
		//	mockBehaviour: func(a adapter) {
		//		testLocation, _ := model.NewLocation(
		//			"test-slug",
		//			"Test LocationDTO for mall",
		//			"Brooklyn, st. main Avenue, 2378",
		//		)
		//		_ = testLocation.Deactivate()
		//		a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocation, nil)
		//		a.location.EXPECT().Update(mock.Anything, mock.AnythingOfType("*model.Location")).Return(nil)
		//	},
		//	expectErr: nil,
		//},
		//{
		//	name: "Success - deactivate location",
		//	input: dto.UpdateLocationInput{
		//		ID:       testID,
		//		Slug:     utils.VPtr("updated-slug"),
		//		Name:     utils.VPtr("Updated Name"),
		//		Address:  utils.VPtr("Updated Address address"),
		//		IsActive: utils.VPtr(false),
		//	},
		//	mockBehaviour: func(a adapter) {
		//		testLocation, _ := model.NewLocation(
		//			"test-slug",
		//			"Test LocationDTO for mall",
		//			"Brooklyn, st. main Avenue, 2378",
		//		)
		//		a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocation, nil)
		//		a.location.EXPECT().Update(mock.Anything, mock.AnythingOfType("*model.Location")).Return(nil)
		//	},
		//	expectErr: nil,
		//},
		{
			name: "Failure - cannot activate location",
			input: dto.UpdateLocationInput{
				ID:       testID,
				Slug:     utils.VPtr("updated-slug"),
				Name:     utils.VPtr("Updated Name"),
				Address:  utils.VPtr("Updated Address address"),
				IsActive: utils.VPtr(true),
			},
			mockBehaviour: func(a adapter) {
				testLocation, _ := model.NewLocation(
					"test-slug",
					"Test LocationDTO for mall",
					"Brooklyn, st. main Avenue, 2378",
				)
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocation, nil)
			},
			expectErr: ucerrs.ErrCannotActivateLocation,
		},
		{
			name: "Failure - cannot deactivate location",
			input: dto.UpdateLocationInput{
				ID:       testID,
				Slug:     utils.VPtr("updated-slug"),
				Name:     utils.VPtr("Updated Name"),
				Address:  utils.VPtr("Updated Address address"),
				IsActive: utils.VPtr(false),
			},
			mockBehaviour: func(a adapter) {
				testLocation, _ := model.NewLocation(
					"test-slug",
					"Test LocationDTO for mall",
					"Brooklyn, st. main Avenue, 2378",
				)
				// Деактивируем локацию, повторная деактивация должна вернуть ошибку
				_ = testLocation.Deactivate()
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocation, nil)
			},
			expectErr: ucerrs.ErrCannotDeactivateLocation,
		},
		{
			name: "Failure - location not found",
			input: dto.UpdateLocationInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(nil, pkgerrs.ErrObjectNotFound)
			},
			expectErr: ucerrs.ErrLocationNotFound,
		},
		{
			name: "Failure - get location db error",
			input: dto.UpdateLocationInput{
				ID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(nil, errors.New("db error"))
			},
			expectErr: ucerrs.ErrGetLocationByIDDB,
		},
		{
			name: "Failure - invalid update input",
			input: dto.UpdateLocationInput{
				ID:   testID,
				Slug: utils.VPtr(""),
				Name: utils.VPtr(""),
			},
			mockBehaviour: func(a adapter) {
				testLocation, _ := model.NewLocation(
					"test-slug",
					"Test LocationDTO for mall",
					"Brooklyn, st. main Avenue, 2378",
				)
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocation, nil)
			},
			expectErr: ucerrs.ErrInvalidInput,
		},
		{
			name: "Failure - update location db error",
			input: dto.UpdateLocationInput{
				ID:      testID,
				Slug:    utils.VPtr("updated-slug"),
				Name:    utils.VPtr("Updated Name"),
				Address: utils.VPtr("Updated Address Address"),
			},
			mockBehaviour: func(a adapter) {
				testLocation, _ := model.NewLocation(
					"test-slug",
					"Test LocationDTO for mall",
					"Brooklyn, st. main Avenue, 2378",
				)
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocation, nil)
				a.location.EXPECT().Update(mock.Anything, mock.AnythingOfType("*model.Location")).Return(errors.New("db error"))
			},
			expectErr: ucerrs.ErrUpdateLocationDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locRepo := mocks.NewMockLocationRepository(t)

			tt.mockBehaviour(adapter{
				location: locRepo,
			})

			uc := usecase.NewUpdateLocationUC(locRepo)

			out, err := uc.Execute(context.Background(), tt.input)

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, out.Location)
				assert.Equal(t, *tt.input.Name, out.Location.Name)
			}
		})
	}
}
