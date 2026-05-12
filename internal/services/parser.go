package services

import (
	"io"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type ParserInterface interface {
	Extract(ioReader io.ReadSeeker) (encryptionMetadata *domain.EncryptionMetadata, err error)
}
