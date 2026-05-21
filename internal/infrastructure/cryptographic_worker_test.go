package infrastructure_test

import (
	"context"
	"crypto/pbkdf2"
	"crypto/sha256"
	"errors"
	"testing"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
	"github.com/UrielJaloto/surgical-rar-recovery/internal/infrastructure"
)

type pbkdf2TestCase struct {
	name        string
	ctx         context.Context
	attempt     string
	wantError   bool
	targetError error
	errMessage  string
}

func TestPbkdf2Worker_TryToBreak(t *testing.T) {
	worker := infrastructure.NewPbkdf2Worker()
	salt := []byte("saltsaltsalt")
	iterations := 32768
	correctPassword := "mypassword123"

	expectedKey, _ := pbkdf2.Key(sha256.New, correctPassword, salt, iterations, 32)

	metadata := &domain.EncryptionMetadata{
		Salt:          salt,
		Iterations:    iterations,
		PasswordCheck: expectedKey,
	}

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []pbkdf2TestCase{
		{
			name:       "Correct password found",
			ctx:        context.Background(),
			attempt:    correctPassword,
			wantError:  true,
			errMessage: "password found",
		},
		{
			name:      "Incorrect password",
			ctx:       context.Background(),
			attempt:   "wrongpassword",
			wantError: false,
		},
		{
			name:        "Context canceled",
			ctx:         canceledCtx,
			attempt:     "anypassword",
			wantError:   true,
			targetError: context.Canceled,
		},
		{
			name:      "Empty password attempt",
			ctx:       context.Background(),
			attempt:   "",
			wantError: false,
		},
	}

	checkTest := func(t *testing.T, testCase pbkdf2TestCase) {
		t.Helper()

		err := worker.TryToBreak(testCase.ctx, testCase.attempt, metadata)

		if (err != nil) != testCase.wantError {
			t.Fatalf("expected error: %t, got: %v", testCase.wantError, err)
		}

		if testCase.targetError != nil && !errors.Is(err, testCase.targetError) {
			t.Errorf("expected error of type %v, got: %v", testCase.targetError, err)
		}

		if testCase.errMessage != "" && err != nil && err.Error() != testCase.errMessage {
			t.Errorf("expected message '%s', got: '%s'", testCase.errMessage, err.Error())
		}
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			checkTest(t, testCase)
		})
	}
}

func BenchmarkPbkdf2Worker_TryToBreak(b *testing.B) {
	worker := infrastructure.NewPbkdf2Worker()
	password := "benchmark_password"
	salt := []byte("salt123456789012")
	iterations := 32768

	expectedKey, _ := pbkdf2.Key(sha256.New, password, salt, iterations, 32)

	metadata := &domain.EncryptionMetadata{
		Salt:          salt,
		Iterations:    iterations,
		PasswordCheck: expectedKey,
	}

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = worker.TryToBreak(ctx, password, metadata)
	}
}
