package infrastructure

import (
	"context"
	"crypto/pbkdf2"
	"crypto/sha256"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type pbkdf2KeyStretcher struct{}

func NewPbkdf2KeyStretcher() *pbkdf2KeyStretcher {
	return &pbkdf2KeyStretcher{}
}

func (p *pbkdf2KeyStretcher) StretchKey(ctx context.Context, passwordAttempt []byte, encryptionMetadata *domain.EncryptionMetadata) (stretchedKey []byte, err error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	//todo: optimize pbkf2 to use
	stretchedKey, err = pbkdf2.Key(sha256.New, string(passwordAttempt), encryptionMetadata.Salt, encryptionMetadata.Iterations, 32)
	if err != nil {
		return make([]byte, 0), err
	}
	return stretchedKey, nil
}
