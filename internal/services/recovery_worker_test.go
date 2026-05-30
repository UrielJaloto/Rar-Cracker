package services_test

import (
	"context"
	"crypto/pbkdf2"
	"crypto/sha256"
	"testing"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
	"github.com/UrielJaloto/surgical-rar-recovery/internal/infrastructure"
	"github.com/UrielJaloto/surgical-rar-recovery/internal/services"
)

type mockKeyStretcher struct {
	expectedPassword string
	fakeHash         []byte
}

func (m *mockKeyStretcher) StretchKey(passwordAttempt []byte, metadata *domain.EncryptionMetadata) ([]byte, error) {
	if string(passwordAttempt) == m.expectedPassword {
		return m.fakeHash, nil
	}
	return []byte("wrong_hash"), nil
}

func TestRecoveryWorker_TryToRecovery(t *testing.T) {
	charset := []rune("abcdefghijklmnopqrstuvwxyz")
	generator := infrastructure.NewPasswordGenerator()

	fakeHash := []byte("success_hash")
	stretcher := &mockKeyStretcher{
		expectedPassword: "basecretaaa",
		fakeHash:         fakeHash,
	}

	metadata := &domain.EncryptionMetadata{
		PasswordCheck: fakeHash,
	}

	config := &domain.Config{
		Charset:   charset,
		KnownPart: "secret",
	}

	chunk := &domain.Chunk{
		VariableLen:    5,
		KnownPartIndex: 2,
		TotalLen:       11,
		StartIndex:     0,
		EndIndex:       1000000,
	}

	worker := services.NewRecoveryWorker(stretcher, generator, metadata, config)
	ctx := context.Background()

	password, err := worker.TryToRecovery(ctx, chunk)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(password) != "basecretaaa" {
		t.Errorf("expected basecretaaa, got %s", string(password))
	}
}

func BenchmarkRecoveryWorker_RealRARScenario(b *testing.B) {
	charset := []rune("abcdefghijklmnopqrstuvwxyz")
	generator := infrastructure.NewPasswordGenerator()
	stretcher := infrastructure.NewPbkdf2KeyStretcher()

	salt := []byte("1234567890123456")
	iterations := 32768
	targetPassword := "zzzz"

	expectedHash, _ := pbkdf2.Key(sha256.New, targetPassword, salt, iterations, 32)

	metadata := &domain.EncryptionMetadata{
		Salt:          salt,
		Iterations:    iterations,
		PasswordCheck: expectedHash,
	}

	config := &domain.Config{
		Charset:   charset,
		KnownPart: "",
	}

	worker := services.NewRecoveryWorker(stretcher, generator, metadata, config)
	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		chunk := &domain.Chunk{
			VariableLen:    4,
			KnownPartIndex: 0,
			TotalLen:       4,
			StartIndex:     uint64(i),
			EndIndex:       uint64(i + 1),
		}

		_, _ = worker.TryToRecovery(ctx, chunk)
	}
}
