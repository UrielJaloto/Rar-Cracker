package services

import (
	"context"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type CryptographicWorkerInterface interface {
	TryToBreak(ctx context.Context, passwordAttempt string, encryptionMetadata *domain.EncryptionMetadata) (err error)
}
