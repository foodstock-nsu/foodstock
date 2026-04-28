package yookassa

type amountDTO struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type confirmationDTO struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type createPaymentRequest struct {
	Amount       amountDTO         `json:"amount"`
	Capture      bool              `json:"capture"`
	Confirmation confirmationDTO   `json:"confirmation"`
	Description  string            `json:"description,omitempty"`
	Metadata     map[string]string `json:"metadata"`
}

type paymentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}
