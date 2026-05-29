package infrastructure

import "github.com/UrielJaloto/surgical-rar-recovery/internal/domain"

type PasswordGenerator struct{}

func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{}
}

func (pg *PasswordGenerator) GeneratePassword(params *domain.PasswordGeneratorParams) (err error) {

	return
}
