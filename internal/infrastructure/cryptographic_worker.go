package infrastructure

import (
	"sync"

	"github.com/UrielJaloto/Rar-Cracker/internal/domain"
	"github.com/UrielJaloto/Rar-Cracker/internal/services"
)

type CryptographicWorker struct {
	passwordAttempt string
	waitGroup       *sync.WaitGroup
	passwordCorrect bool
}

func NewCryptographicWorker(passwordAttempt string, waitgroup *sync.WaitGroup) services.CryptographicWorkerInterface {
	return &CryptographicWorker{passwordAttempt, waitgroup, false}
}

func (c *CryptographicWorker) DeriveKey(config *domain.EncryptionMetadata) (err error) {
	return nil
	//todo
}
