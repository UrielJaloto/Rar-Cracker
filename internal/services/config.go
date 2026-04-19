package services

import "github.com/UrielJaloto/Rar-Cracker/internal/domain"

type ConfigLoaderInterface interface {
	Load() (*domain.Config, error)
}

type ConfigValidatorInterface interface {
	Validate(settings *domain.Config) ([]string, error)
}
