package services

import "context"

type CryptographicWorkerInterface interface {
	TryToBreak(ctx context.Context, passwordAttempt string) (err error)
}
