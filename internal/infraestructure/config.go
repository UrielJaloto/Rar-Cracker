package infraestructure

import "github.com/UrielJaloto/Rar-Cracker/internal/domain"

type ConfigInterface interface {
	Load() (*domain.Config, error)
}
