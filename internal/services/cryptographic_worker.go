package services

import "github.com/UrielJaloto/Rar-Cracker/internal/domain"

type CryptographicWorkerInterface interface {
	DeriveKey(config *domain.EncryptionMetadata) (err error)
}
