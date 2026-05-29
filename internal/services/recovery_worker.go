package services

import (
	"context"
	"errors"
	"slices"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type RecoveryWorker struct {
	keyStretcher      KeyStretcherInterface
	passwordGenerator PasswordGeneratorInterface
}

func NewRecoveryWorker(keyStretcher KeyStretcherInterface, passwordGenerator PasswordGeneratorInterface) *RecoveryWorker {
	return &RecoveryWorker{keyStretcher, passwordGenerator}
}

func (rw *RecoveryWorker) TryToRecovery(ctx context.Context, encryptionMetadata *domain.EncryptionMetadata, config *domain.Config, chunk *domain.Chunk) (password []byte, err error) {
	var stretchedKey []byte

	buffer := make([]byte, chunk.TotalLen)
	copy(buffer[chunk.KnownPartIndex:], []byte(config.KnownPart))

	for iteration := chunk.StartIndex; iteration < chunk.EndIndex; iteration++ {
		if err = ctx.Err(); err != nil {
			return password, err
		}

		PasswordGeneratorParams := domain.PasswordGeneratorParams{
			Buffer:    buffer,
			Charset:   config.Charset,
			Iteration: iteration,
			Chunk:     chunk,
		}

		rw.passwordGenerator.GeneratePassword(&PasswordGeneratorParams)
		passwordAttempt := buffer[:chunk.TotalLen]

		if stretchedKey, err = rw.keyStretcher.StretchKey(passwordAttempt, encryptionMetadata); err != nil {
			return password, err
		}

		if slices.Equal(stretchedKey, encryptionMetadata.PasswordCheck) {
			password = make([]byte, chunk.TotalLen)
			copy(password, passwordAttempt)
			return password, err
		}
	}

	err = errors.New("wrong password")
	return password, err
}
