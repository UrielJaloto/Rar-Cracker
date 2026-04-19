package services

import (
	"io"

	"github.com/UrielJaloto/Rar-Cracker/internal/domain"
)

type ParserInterface interface {
	Extract(ioReader io.ReadSeeker) (*domain.EncryptionMetadata, error)
}
