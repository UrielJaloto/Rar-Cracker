package services

import "context"

type PasswordGeneratorInterface interface {
	GeneratePassword(ctx context.Context, iteration uint64, variableLen int, KnownPartIndex int, buffer []byte) (err error)
}
