package port

type QRCodeGenerator interface {
	Generate(slug string) ([]byte, error)
}
