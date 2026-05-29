package infrastructure_test

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"errors"
	"slices"
	"testing"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
	"github.com/UrielJaloto/surgical-rar-recovery/internal/infrastructure"
)

type pbkdf2TestCase struct {
	name        string
	attempt     []byte
	wantError   bool
	wantEqual   bool
	targetError error
	errMessage  string
}

func TestPbkdf2Worker_StretchKey(t *testing.T) {
	service := infrastructure.NewPbkdf2KeyStretcher()
	salt := []byte("saltsaltsalt")
	iterations := 32768
	correctPassword := []byte("mypassword123")

	expectedKey, _ := pbkdf2.Key(sha256.New, string(correctPassword), salt, iterations, 32)

	metadata := &domain.EncryptionMetadata{
		Salt:          salt,
		Iterations:    iterations,
		PasswordCheck: expectedKey,
	}

	tests := []pbkdf2TestCase{
		{
			name:       "Correct password found",
			attempt:    correctPassword,
			wantError:  false,
			errMessage: "password found",
			wantEqual:  true,
		},
		{
			name:      "Incorrect password",
			attempt:   []byte("wrongpassword"),
			wantError: false,
			wantEqual: false,
		},
		{
			name:      "Empty password attempt",
			attempt:   []byte(""),
			wantError: false,
			wantEqual: false,
		},
	}

	checkTest := func(t *testing.T, testCase pbkdf2TestCase) {
		t.Helper()

		stretchedKey, err := service.StretchKey(testCase.attempt, metadata)

		isEqual := slices.Equal(stretchedKey, metadata.PasswordCheck)

		if (err != nil) != testCase.wantError {
			t.Fatalf("expected error: %t, got: %v", testCase.wantError, err)
		}

		if (isEqual != false) != testCase.wantEqual {
			t.Fatalf("expected equal: %t, got: %v", testCase.wantEqual, isEqual)
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

func BenchmarkPbkdf2Worker_StretchKey(b *testing.B) {
	service := infrastructure.NewPbkdf2KeyStretcher()
	password := []byte("benchmark_password")
	salt := []byte("salt123456789012")
	iterations := 32768

	expectedKey, _ := pbkdf2.Key(sha256.New, string(password), salt, iterations, 32)

	metadata := &domain.EncryptionMetadata{
		Salt:          salt,
		Iterations:    iterations,
		PasswordCheck: expectedKey,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = service.StretchKey(password, metadata)
	}
}
