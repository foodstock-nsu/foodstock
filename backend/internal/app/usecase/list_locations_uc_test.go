package usecase_test

import (
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/usecase"
	"backend/internal/domain/model"
	"backend/internal/domain/port/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListLocationsUC_Execute(t *testing.T) {
	type adapter struct {
		location *mocks.MockLocationRepository
	}

	type testCase struct {
		name          string
		mockBehaviour func(a adapter)
		expectErr     error
	}

	testLocation, _ := model.NewLocation(
		"test-slug",
		"Test Location for mall",
		"Brooklyn, st. main Avenue, 2378",
	)

	var tests = []testCase{
		{
			name: "Success",
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().List(mock.Anything).Return([]*model.Location{testLocation}, nil)
			},
			expectErr: nil,
		},
		{
			name: "Failure - list locations db error",
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().List(mock.Anything).Return(nil, errors.New("db error"))
			},
			expectErr: ucerrs.ErrListLocationsDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locRepo := mocks.NewMockLocationRepository(t)

			tt.mockBehaviour(adapter{
				location: locRepo,
			})

			uc := usecase.NewListLocationsUC(locRepo)

			out, err := uc.Execute(context.Background())

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, out.Locations)
				assert.Equal(t, testLocation.Name(), out.Locations[0].Name)
			}
		})
	}
}
