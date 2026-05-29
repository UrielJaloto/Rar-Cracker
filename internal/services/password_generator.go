package services

import (
	"context"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type PasswordGeneratorInterface interface {
	GeneratePassword(ctx context.Context, buffer []byte, chunk *domain.Chunk, iteration uint64) (err error)
}
