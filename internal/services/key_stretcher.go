package services

import (
	"context"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type KeyStretcherInterface interface {
	StretchKey(ctx context.Context, passwordAttempt []byte, encryptionMetadata *domain.EncryptionMetadata) (stretchedKey []byte, err error)
}
