package services

import (
	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type KeyStretcherInterface interface {
	StretchKey(passwordAttempt []byte, encryptionMetadata *domain.EncryptionMetadata) (stretchedKey []byte, err error)
}
