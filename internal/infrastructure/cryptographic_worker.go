package infrastructure

import (
	"context"
	"crypto/pbkdf2"
	"crypto/sha256"
	"errors"
	"slices"
	"sync"

	"github.com/UrielJaloto/Rar-Cracker/internal/domain"
	"github.com/UrielJaloto/Rar-Cracker/internal/services"
)

type pbkdf2Worker struct {
	EncryptionMetadata *domain.EncryptionMetadata
	waitGroup          *sync.WaitGroup
	stretchedKey       []byte
}

func NewPbkdf2Worker(encryptionMetadata *domain.EncryptionMetadata, waitgroup *sync.WaitGroup) services.CryptographicWorkerInterface {
	return &pbkdf2Worker{encryptionMetadata, waitgroup, make([]byte, 0)}
}

func (p *pbkdf2Worker) TryToBreak(ctx context.Context, passwordAttempt string) (err error) {
	defer p.waitGroup.Done()

	if err := ctx.Err(); err != nil {
		return err
	}

	if err = p.StretchKey(ctx, passwordAttempt); err != nil {
		return err
	}

	if slices.Equal(p.stretchedKey, p.EncryptionMetadata.PasswordCheck) {
		return errors.New("password found")
	}

	return nil
}

func (p *pbkdf2Worker) StretchKey(ctx context.Context, passwordAttempt string) (err error) {
	p.stretchedKey, err = pbkdf2.Key(sha256.New, passwordAttempt, p.EncryptionMetadata.Salt, p.EncryptionMetadata.Iterations, 32)
	if err != nil {
		return err
	}
	return nil
}
