package dto

type GetQRCodeInput struct {
	Slug string
}

type GetQRCodeOutput struct {
	QRCode []byte
}
