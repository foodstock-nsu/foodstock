package qrcode

import (
	"fmt"

	goqrcode "github.com/skip2/go-qrcode"
)

type Generator struct {
	baseURL string
	size    int // size of qr-code image
}

func NewGenerator(baseURL string, size int) *Generator {
	if size <= 0 {
		size = 256
	}
	return &Generator{
		baseURL: baseURL,
		size:    size,
	}
}

func (g *Generator) Generate(slug string) ([]byte, error) {
	url := fmt.Sprintf("%s/l/%s", g.baseURL, slug)

	png, err := goqrcode.Encode(url, goqrcode.Medium, g.size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate qr code: %w", err)
	}

	return png, nil
}
