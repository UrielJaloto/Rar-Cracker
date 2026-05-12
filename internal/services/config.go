package services

import "github.com/UrielJaloto/surgical-rar-recovery/internal/domain"

type ConfigLoaderInterface interface {
	Load() (settings *domain.Config, err error)
}

type ConfigValidatorInterface interface {
	Validate(settings *domain.Config) (validationReport domain.ValidationReport)
}
