package infrastructure

import (
	"context"
	"crypto/pbkdf2"
	"crypto/sha256"
	"errors"
	"slices"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type pbkdf2Worker struct{}

func NewPbkdf2Worker() *pbkdf2Worker {
	return &pbkdf2Worker{}
}

func (p *pbkdf2Worker) TryToBreak(ctx context.Context, passwordAttempt string, encryptionMetadata *domain.EncryptionMetadata) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	stretchedKey, err := p.StretchKey(ctx, passwordAttempt, encryptionMetadata)
	if err != nil {
		return err
	}

	if slices.Equal(stretchedKey, encryptionMetadata.PasswordCheck) {
		return errors.New("password found")
	}

	return nil
}

func (p *pbkdf2Worker) StretchKey(ctx context.Context, passwordAttempt string, encryptionMetadata *domain.EncryptionMetadata) (stretchedKey []byte, err error) {
	stretchedKey, err = pbkdf2.Key(sha256.New, passwordAttempt, encryptionMetadata.Salt, encryptionMetadata.Iterations, 32)
	if err != nil {
		return make([]byte, 0), err
	}
	return stretchedKey, nil
}
