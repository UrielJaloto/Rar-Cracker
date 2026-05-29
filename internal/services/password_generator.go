package services

import (
	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type PasswordGeneratorInterface interface {
	GeneratePassword(params *domain.PasswordGeneratorParams) (err error)
}
