package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"backend/internal/adapter/out/postgres"
	"backend/internal/domain/model"
	"backend/internal/infrastructure/password" // Твой BcryptHasher

	"github.com/google/uuid"
)

// syncAdmins — техническая функция инициализации
func syncAdmins(ctx context.Context, repo *postgres.AdminRepository, password *password.Hasher) error {
	// Достаем строку вида "admin:pass123,manager:pass456"
	rawConfig := os.Getenv("ADMIN_SETUP")
	if rawConfig == "" {
		return nil
	}

	entries := strings.Split(rawConfig, ",")
	for _, entry := range entries {
		parts := strings.Split(entry, ":")
		if len(parts) != 2 {
			continue
		}

		login, pass := parts[0], parts[1]

		// Хешируем
		hash, err := password.Hash(pass)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", login, err)
		}

		// Создаем доменную модель
		// UUID берем новый, т.к. Upsert в БД сориентируется по Login
		admin := model.RestoreAdmin(uuid.New(), login, hash, time.Now())

		if err := repo.Upsert(ctx, admin); err != nil {
			return fmt.Errorf("failed to upsert admin %s: %w", login, err)
		}
	}

	return nil
}
