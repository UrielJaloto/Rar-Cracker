package infrastructure

import (
	"context"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type PasswordGenerator struct{}

func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{}
}

func (pg *PasswordGenerator) GeneratePassword(ctx context.Context, buffer []byte, chunk *domain.Chunk, iteration uint64) (err error) {
	if err = ctx.Err(); err != nil {
		return
	}

	// var v = chunk.VariableLen
	// var x = chunk.KnownPartIndex

	//todo
	return
}
