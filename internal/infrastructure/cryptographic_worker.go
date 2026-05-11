package infrastructure

import (
	"context"
	"crypto/pbkdf2"
	"crypto/sha256"
	"errors"
	"slices"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
	"github.com/UrielJaloto/surgical-rar-recovery/internal/services"
)

type pbkdf2Worker struct {
	EncryptionMetadata *domain.EncryptionMetadata
}

func NewPbkdf2Worker(encryptionMetadata *domain.EncryptionMetadata) services.CryptographicWorkerInterface {
	return &pbkdf2Worker{encryptionMetadata}
}

func (p *pbkdf2Worker) TryToBreak(ctx context.Context, passwordAttempt string) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	stretchedKey, err := p.StretchKey(ctx, passwordAttempt)
	if err != nil {
		return err
	}

	if slices.Equal(stretchedKey, p.EncryptionMetadata.PasswordCheck) {
		return errors.New("password found")
	}

	return nil
}

func (p *pbkdf2Worker) StretchKey(ctx context.Context, passwordAttempt string) (stretchedKey []byte, err error) {
	stretchedKey, err = pbkdf2.Key(sha256.New, passwordAttempt, p.EncryptionMetadata.Salt, p.EncryptionMetadata.Iterations, 32)
	if err != nil {
		return make([]byte, 0), err
	}
	return stretchedKey, nil
}
