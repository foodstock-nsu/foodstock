package usecase_test

import (
	"backend/internal/app/dto"
	ucerrs "backend/internal/app/errs"
	"backend/internal/app/usecase"
	"backend/internal/domain/model"
	"backend/internal/domain/port/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCatalogUC_Execute(t *testing.T) {
	type adapter struct {
		item         *mocks.MockItemRepository
		locationItem *mocks.MockLocationItemRepository
	}

	type testCase struct {
		name          string
		input         dto.GetCatalogInput
		mockBehaviour func(a adapter)
		expectOutput  dto.GetCatalogOutput
		expectErr     error
	}

	testLocationID := uuid.New()
	testItemID1 := uuid.New()
	testItemID2 := uuid.New()

	var tests = []testCase{
		{
			name: "Success - catalog with multiple items and categories",
			input: dto.GetCatalogInput{
				LocationID: testLocationID,
			},
			mockBehaviour: func(a adapter) {
				locationItems := []*model.LocationItem{
					model.RestoreLocationItem(
						uuid.New(),
						testItemID1,
						testLocationID,
						0,
						false,
						0,
					),
					model.RestoreLocationItem(
						uuid.New(),
						testItemID2,
						testLocationID,
						0,
						false,
						0,
					),
				}

				items := []*model.Item{
					model.RestoreItem(
						testItemID1,
						"Super sandwich",
						nil,
						model.ItemLunch,
						"https://photo.jpg",
						nil,
						time.Now().UTC(),
					),
					model.RestoreItem(
						testItemID2,
						"Super chicken roll",
						nil,
						model.ItemBreakfast,
						"https://photo.jpg",
						nil,
						time.Now().UTC(),
					),
				}

				a.locationItem.EXPECT().List(mock.Anything, testLocationID).Return(locationItems, nil)
				a.item.EXPECT().ListByIDs(mock.Anything, []uuid.UUID{testItemID1, testItemID2}).Return(items, nil)
			},
			expectOutput: dto.GetCatalogOutput{
				Categories: []string{"Lunch", "Breakfast"},
				Items:      []dto.CatalogItemOutput{{}, {}},
			},
			expectErr: nil,
		},
		{
			name:  "Success - catalog with items from same category",
			input: dto.GetCatalogInput{LocationID: testLocationID},
			mockBehaviour: func(a adapter) {
				locationItems := []*model.LocationItem{
					model.RestoreLocationItem(
						uuid.New(),
						testItemID1,
						testLocationID,
						0,
						false,
						0,
					),
					model.RestoreLocationItem(
						uuid.New(),
						testItemID2,
						testLocationID,
						0,
						false,
						0,
					),
				}

				items := []*model.Item{
					model.RestoreItem(
						testItemID1,
						"Super sandwich",
						nil,
						model.ItemLunch,
						"https://photo.jpg",
						nil,
						time.Now().UTC(),
					),
					model.RestoreItem(
						testItemID2,
						"Super chicken roll",
						nil,
						model.ItemLunch,
						"https://photo.jpg",
						nil,
						time.Now().UTC(),
					),
				}

				a.locationItem.EXPECT().List(mock.Anything, testLocationID).Return(locationItems, nil)
				a.item.EXPECT().ListByIDs(mock.Anything, []uuid.UUID{testItemID1, testItemID2}).Return(items, nil)
			},
			expectOutput: dto.GetCatalogOutput{
				Categories: []string{"Lunch"},
				Items:      []dto.CatalogItemOutput{{}, {}},
			},
			expectErr: nil,
		},
		{
			name:  "Success - empty inventory",
			input: dto.GetCatalogInput{LocationID: testLocationID},
			mockBehaviour: func(a adapter) {
				a.locationItem.EXPECT().List(mock.Anything, testLocationID).Return([]*model.LocationItem{}, nil)
			},
			expectOutput: dto.GetCatalogOutput{},
			expectErr:    nil,
		},
		{
			name:  "Failure - list location items db error",
			input: dto.GetCatalogInput{LocationID: testLocationID},
			mockBehaviour: func(a adapter) {
				a.locationItem.EXPECT().List(mock.Anything, testLocationID).Return(nil, errors.New("db error"))
			},
			expectOutput: dto.GetCatalogOutput{},
			expectErr:    ucerrs.ErrListLocationItemsDB,
		},
		{
			name:  "Failure - list items by IDs db error",
			input: dto.GetCatalogInput{LocationID: testLocationID},
			mockBehaviour: func(a adapter) {
				locationItems := []*model.LocationItem{
					model.RestoreLocationItem(
						uuid.New(),
						testItemID1,
						testLocationID,
						0,
						false,
						0,
					),
				}

				a.locationItem.EXPECT().List(mock.Anything, testLocationID).Return(locationItems, nil)
				a.item.EXPECT().ListByIDs(mock.Anything, []uuid.UUID{testItemID1}).Return(nil, errors.New("db error"))
			},
			expectOutput: dto.GetCatalogOutput{},
			expectErr:    ucerrs.ErrListItemsByIDsDB,
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

			uc := usecase.NewGetCatalogUC(itemRepo, locItemRepo)

			output, err := uc.Execute(context.Background(), tt.input)

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Equal(t, tt.expectOutput, output)
			} else {
				assert.NoError(t, err)
				assert.Len(t, output.Categories, len(tt.expectOutput.Categories))
				assert.Len(t, output.Items, len(tt.expectOutput.Items))
			}
		})
	}
}
