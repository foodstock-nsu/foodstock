package http

import (
	"backend/internal/domain/port"
)

type Router struct {
	tokenGen port.TokenGenerator
}
