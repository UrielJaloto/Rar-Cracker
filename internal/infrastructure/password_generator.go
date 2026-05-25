package infrastructure

import (
	"context"
)

type PasswordGenerator struct{}

func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{}
}

func (pg *PasswordGenerator) GeneratePassword(ctx context.Context, iteration uint64, variableLen int, KnownPartIndex int, buffer []byte) (err error) {
	if err = ctx.Err(); err != nil {
		return
	}

	//todo
	return
}
