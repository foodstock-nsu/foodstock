package model_test

import (
	"backend/internal/domain/model"
	pkgerrs "backend/pkg/errs"
	"backend/pkg/utils"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNutrition(t *testing.T) {
	var (
		testCalories = utils.VPtr(1000)
		testProteins = utils.VPtr(float64(10))
		testFats     = utils.VPtr(float64(10))
		testCarbs    = utils.VPtr(float64(10))
	)

	type testCase struct {
		testName string
		calories *int
		proteins *float64
		fats     *float64
		carbs    *float64
		expect   error
	}

	var testCases = []testCase{
		{
			testName: "Success",
			calories: testCalories,
			proteins: testProteins,
			fats:     testFats,
			carbs:    testCarbs,
			expect:   nil,
		},
		{
			testName: "Success - empty nutrition",
			expect:   nil,
		},
		{
			testName: "Failure - invalid calories",
			calories: utils.VPtr(-1),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid proteins",
			calories: testCalories,
			proteins: utils.VPtr(float64(-1)),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid fats",
			calories: testCalories,
			proteins: testProteins,
			fats:     utils.VPtr(float64(-1)),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid carbs",
			calories: testCalories,
			proteins: testProteins,
			fats:     testFats,
			carbs:    utils.VPtr(float64(-1)),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			nutrition, err := model.NewNutrition(
				tt.calories, tt.proteins, tt.fats, tt.carbs,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Empty(t, nutrition)
			} else {
				require.NoError(t, err)

				if tt.calories == nil && tt.proteins == nil && tt.fats == nil && tt.carbs == nil {
					assert.Nil(t, nutrition)
				} else {
					assert.Equal(t, tt.calories, nutrition.Calories())
					assert.Equal(t, tt.proteins, nutrition.Proteins())
					assert.Equal(t, tt.fats, nutrition.Fats())
					assert.Equal(t, tt.carbs, nutrition.Carbs())
				}
			}
		})
	}
}

func TestNutrition_Update(t *testing.T) {
	var (
		testCalories = utils.VPtr(1000)
		testProteins = utils.VPtr(float64(10))
		testFats     = utils.VPtr(float64(10))
		testCarbs    = utils.VPtr(float64(10))
	)

	type testCase struct {
		testName string
		calories *int
		proteins *float64
		fats     *float64
		carbs    *float64
		expect   error
	}

	var testCases = []testCase{
		{
			testName: "Success",
			calories: testCalories,
			proteins: testProteins,
			fats:     testFats,
			carbs:    testCarbs,
			expect:   nil,
		},
		{
			testName: "Failure - invalid calories",
			calories: utils.VPtr(-10),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid proteins",
			proteins: utils.VPtr(float64(-1)),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid fats",
			fats:     utils.VPtr(float64(-1)),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid carbs",
			carbs:    utils.VPtr(float64(-1)),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			nutrition := model.RestoreNutrition(
				utils.VPtr(1000),
				utils.VPtr(float64(10)),
				utils.VPtr(float64(10)),
				utils.VPtr(float64(10)),
			)

			err := nutrition.Update(
				tt.calories, tt.proteins, tt.fats, tt.carbs,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				if tt.calories != nil {
					assert.NotEqual(t, tt.calories, nutrition.Calories())
				}
				if tt.proteins != nil {
					assert.NotEqual(t, tt.proteins, nutrition.Proteins())
				}
				if tt.fats != nil {
					assert.NotEqual(t, tt.fats, nutrition.Fats())
				}
				if tt.carbs != nil {
					assert.NotEqual(t, tt.carbs, nutrition.Carbs())
				}
			} else {
				require.NoError(t, err)
				if tt.calories != nil {
					assert.Equal(t, tt.calories, nutrition.Calories())
					fmt.Printf("\nIn tt = %v | in model = %v\n", *tt.calories, *nutrition.Calories())
				}
				if tt.proteins != nil {
					assert.Equal(t, tt.proteins, nutrition.Proteins())
				}
				if tt.fats != nil {
					assert.Equal(t, tt.fats, nutrition.Fats())
				}
				if tt.carbs != nil {
					assert.Equal(t, tt.carbs, nutrition.Carbs())
				}
			}
		})
	}
}

func TestNewItem(t *testing.T) {
	var (
		testItemName  = "Chicken Sandwich"
		testDesc      = utils.VPtr("Chicken sandwich with fresh vegetables")
		testCategory  = "lunch"
		testPhoto     = "https://hosting.com/new.jpg"
		testNutrition = model.RestoreNutrition(
			utils.VPtr(1000),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(10)),
		)
	)

	type testCase struct {
		testName    string
		name        string
		description *string
		category    string
		photoURL    string
		nutrition   *model.Nutrition
		expect      error
	}

	var testCases = []testCase{
		{
			testName:    "Success",
			name:        testItemName,
			description: testDesc,
			category:    testCategory,
			photoURL:    testPhoto,
			nutrition:   testNutrition,
			expect:      nil,
		},
		{
			testName: "Failure - invalid name",
			name:     "inv",
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName:    "Failure - invalid description",
			name:        testItemName,
			description: utils.VPtr("new desc"), // too short
			expect:      pkgerrs.ErrValueIsInvalid,
		},
		{
			testName:    "Failure - invalid category",
			name:        testItemName,
			description: testDesc,
			category:    "unexisting",
			expect:      pkgerrs.ErrValueIsInvalid,
		},
		{
			testName:    "Failure - invalid photo url",
			name:        testItemName,
			description: testDesc,
			category:    testCategory,
			photoURL:    "photo.jpg",
			expect:      pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			item, err := model.NewItem(
				tt.name,
				tt.description,
				tt.category,
				tt.photoURL,
				tt.nutrition,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
				assert.Nil(t, item)
			} else {
				require.NoError(t, err)
				require.NotNil(t, item)

				assert.NotEmpty(t, item.ID())
				assert.Equal(t, tt.name, item.Name())
				if tt.description != nil {
					assert.Equal(t, tt.description, item.Description())
				}
				assert.Equal(t, tt.category, item.Category().String())
				assert.Equal(t, tt.photoURL, item.PhotoURL())
				if tt.nutrition != nil {
					assert.Equal(t, tt.nutrition, item.Nutrition())
				}
				assert.False(t, item.CreatedAt().After(time.Now().UTC()))
			}
		})
	}
}

func TestItem_Update(t *testing.T) {
	var (
		testItemName  = utils.VPtr("Chicken Sandwich")
		testDesc      = utils.VPtr("Chicken sandwich with fresh vegetables")
		testCategory  = utils.VPtr("lunch")
		testPhoto     = utils.VPtr("https://hosting.com/new.jpg")
		testNutrition = model.RestoreNutrition(
			utils.VPtr(1000),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(10)),
			utils.VPtr(float64(10)),
		)
	)

	type testCase struct {
		testName    string
		name        *string
		description *string
		category    *string
		photoURL    *string
		nutrition   *model.Nutrition
		expect      error
	}

	var testCases = []testCase{
		{
			testName:    "Success",
			name:        testItemName,
			description: testDesc,
			category:    testCategory,
			photoURL:    testPhoto,
			nutrition:   testNutrition,
			expect:      nil,
		},
		{
			testName: "Failure - invalid name",
			name:     utils.VPtr("inv"),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName:    "Failure - invalid description",
			description: utils.VPtr("new desc"),
			expect:      pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid category",
			category: utils.VPtr("unexisting"),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
		{
			testName: "Failure - invalid photo url",
			photoURL: utils.VPtr("photo.jpg"),
			expect:   pkgerrs.ErrValueIsInvalid,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			item := model.RestoreItem(
				uuid.New(),
				"Salade",
				nil,
				model.ItemBreakfast,
				"https://hosting.com/old.jpg",
				model.RestoreNutrition(
					utils.VPtr(1000),
					utils.VPtr(float64(10)),
					utils.VPtr(float64(10)),
					utils.VPtr(float64(10)),
				),
				time.Now().UTC(),
			)

			err := item.Update(
				tt.name, tt.description, tt.category,
				tt.photoURL, tt.nutrition,
			)
			if tt.expect != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expect)
			} else {
				require.NoError(t, err)
				if tt.name != nil {
					assert.Equal(t, *tt.name, item.Name())
				}
				if tt.description != nil {
					assert.Equal(t, tt.description, item.Description())
				}
				if tt.category != nil {
					assert.Equal(t, *tt.category, item.Category().String())
				}
				if tt.photoURL != nil {
					assert.Equal(t, *tt.photoURL, item.PhotoURL())
				}
				if tt.nutrition != nil {
					assert.Equal(t, tt.nutrition, item.Nutrition())
				}
			}
		})
	}
}
