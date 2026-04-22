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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetQRCodeUC_Execute(t *testing.T) {
	type adapter struct {
		location *mocks.MockLocationRepository
		qrcode   *mocks.MockQRCodeGenerator
	}

	type testCase struct {
		name          string
		input         dto.GetQRCodeInput
		mockBehaviour func(a adapter)
		expectErr     error
	}

	testID := uuid.New()
	testLocationOperational, _ := model.NewLocation(
		"test-slug",
		"Test LocationOutput for mall",
		"Brooklyn, st. main Avenue, 2378",
	)
	testLocationNonOperational := &model.Location{}

	var tests = []testCase{
		{
			name: "Success",
			input: dto.GetQRCodeInput{
				LocationID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocationOperational, nil)
				a.qrcode.EXPECT().Generate(mock.Anything).Return([]byte("qr-code-data"), nil)
			},
			expectErr: nil,
		},
		{
			name: "Failure - db error",
			input: dto.GetQRCodeInput{
				LocationID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(nil, errors.New("db error"))
			},
			expectErr: ucerrs.ErrGetLocationByIDDB,
		},
		{
			name: "Failure - not operational",
			input: dto.GetQRCodeInput{
				LocationID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocationNonOperational, nil)
			},
			expectErr: ucerrs.ErrCannotGetLocationQRCode,
		},
		{
			name: "Failure - qr generation error",
			input: dto.GetQRCodeInput{
				LocationID: testID,
			},
			mockBehaviour: func(a adapter) {
				a.location.EXPECT().GetByID(mock.Anything, testID).Return(testLocationOperational, nil)
				a.qrcode.EXPECT().Generate(mock.Anything).Return(nil, errors.New("qr error"))
			},
			expectErr: ucerrs.ErrGenerateQRCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locRepo := mocks.NewMockLocationRepository(t)
			qrGen := mocks.NewMockQRCodeGenerator(t)

			tt.mockBehaviour(adapter{
				location: locRepo,
				qrcode:   qrGen,
			})

			uc := usecase.NewGetQRCodeUC(locRepo, qrGen)

			out, err := uc.Execute(context.Background(), tt.input)

			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, []byte("qr-code-data"), out.QRCode)
			}
		})
	}
}
