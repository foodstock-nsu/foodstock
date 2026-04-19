package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const leewayVal = 5 * time.Second

type TokenGenerator struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenGenerator(secret string, ttl time.Duration) *TokenGenerator {
	return &TokenGenerator{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (gen *TokenGenerator) Generate(adminID uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   adminID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(gen.ttl).UTC()),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(gen.secret)
}

func (gen *TokenGenerator) parseToken(token string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return gen.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), jwt.WithLeeway(leewayVal))

	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return claims, nil
}

func (gen *TokenGenerator) Validate(token string) (uuid.UUID, error) {
	claims, err := gen.parseToken(token)
	if err != nil {
		return uuid.Nil, err
	}

	adminID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get admin id: %w", err)
	}

	return adminID, nil
}
