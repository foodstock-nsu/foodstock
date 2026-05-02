package yookassa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type PaymentGateway struct {
	shopID     string
	apiKey     string
	httpClient *http.Client
}

func NewPaymentGateway(shopID, apiKey string, timeout time.Duration) *PaymentGateway {
	return &PaymentGateway{
		shopID: shopID,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

var (
	ErrAuthInvalid     = errors.New("failed to fetch yookassa api: invalid shop id or api key")
	ErrPaymentNotFound = errors.New("failed to fetch yookassa api: payment not found")
	ErrTooManyRequests = errors.New("failed to fetch yookassa api: too many requests")
)

func (p *PaymentGateway) parseYookassaError(errCode int, originalErr error) error {
	msg := fmt.Sprintf("status code %d", errCode)
	ctxErr := errors.New(msg)
	if originalErr != nil {
		ctxErr = originalErr
	}

	switch errCode {
	case 401:
		return ErrAuthInvalid
	case 404:
		return ErrPaymentNotFound
	case 429:
		return ErrTooManyRequests
	default:
		return fmt.Errorf("failed to fetch yookassa api (status %d): %v", errCode, ctxErr)
	}
}

// Create Creates a Payment object and returns its external id and payment url
func (p *PaymentGateway) Create(
	ctx context.Context, amount int64,
	returnURL string, orderID uuid.UUID,
) (string, string, error) {
	const apiURL = "https://api.yookassa.ru/v3/payments"

	amountStr := fmt.Sprintf("%.2f", float64(amount)/100)

	body := createPaymentRequest{
		Amount: amountDTO{
			Value:    amountStr,
			Currency: "RUB",
		},
		Capture: true,
		Confirmation: confirmationDTO{
			Type:      "redirect",
			ReturnURL: returnURL,
		},
		Description: "Оплата корзины",
		Metadata: map[string]string{
			"order_id": orderID.String(),
		},
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodPost, apiURL, bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.SetBasicAuth(p.shopID, p.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	httpReq.Header.Set("Idempotence-Key", uuid.New().String())

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch yookassa api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", p.parseYookassaError(resp.StatusCode, err)
	}

	var result paymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.ID, result.Confirmation.ConfirmationURL, nil
}

func (p *PaymentGateway) GetStatus(ctx context.Context, externalID string) (string, error) {
	apiURL := fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s", externalID)

	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodGet, apiURL, nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.SetBasicAuth(p.shopID, p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to fetch yookassa api: %w", err)
	}
	defer resp.Body.Close()

	var result paymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Status, nil
}
