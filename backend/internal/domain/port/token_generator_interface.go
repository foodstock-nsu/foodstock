package port

import "github.com/google/uuid"

type TokenGenerator interface {
	Generate(adminID uuid.UUID) (string, error)
	Validate(token string) (uuid.UUID, error)
}
