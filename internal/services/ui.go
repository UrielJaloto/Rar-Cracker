package services

import "github.com/UrielJaloto/Rar-Cracker/internal/domain"

type UiInterface interface {
	ShowWarnings(warnings []string)
	ShowErrors(err error)
	ShowConfiguration(config *domain.Config)
	ShowEncryptionMetadata(metaData *domain.EncryptionMetadata)
}
