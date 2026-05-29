package infrastructure_test

import (
	"testing"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
	"github.com/UrielJaloto/surgical-rar-recovery/internal/infrastructure"
)

func TestPasswordGenerator_GeneratePassword(t *testing.T) {
	charset := []rune("abcdefghijklmnopqrstuvwxyz")
	generator := infrastructure.NewPasswordGenerator()

	tests := []struct {
		name           string
		iteration      uint64
		variableLen    int
		totalLen       int
		knownPartIndex int
		knownPart      string
		expected       string
	}{
		{
			name:           "Without known part",
			iteration:      0,
			variableLen:    3,
			totalLen:       3,
			knownPartIndex: 0,
			knownPart:      "",
			expected:       "aaa",
		},
		{
			name:           "Known part at the beginning",
			iteration:      1,
			variableLen:    2,
			totalLen:       6,
			knownPartIndex: 0,
			knownPart:      "test",
			expected:       "testba",
		},
		{
			name:           "Known part at the end",
			iteration:      2,
			variableLen:    3,
			totalLen:       7,
			knownPartIndex: 3,
			knownPart:      "2024",
			expected:       "caa2024",
		},
		{
			name:           "Known part in the middle",
			iteration:      27,
			variableLen:    3,
			totalLen:       6,
			knownPartIndex: 2,
			knownPart:      "mid",
			expected:       "bbmida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := make([]byte, tt.totalLen)

			if tt.knownPart != "" {
				copy(buffer[tt.knownPartIndex:], []byte(tt.knownPart))
			}

			chunk := &domain.Chunk{
				VariableLen:    tt.variableLen,
				KnownPartIndex: tt.knownPartIndex,
				TotalLen:       tt.totalLen,
			}

			params := &domain.PasswordGeneratorParams{
				Buffer:    buffer,
				Charset:   charset,
				Iteration: tt.iteration,
				Chunk:     chunk,
			}

			generator.GeneratePassword(params)

			result := string(buffer)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func BenchmarkPasswordGenerator_GeneratePassword(b *testing.B) {
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	generator := infrastructure.NewPasswordGenerator()

	buffer := make([]byte, 10)
	copy(buffer[4:], []byte("known"))

	chunk := &domain.Chunk{
		VariableLen:    5,
		KnownPartIndex: 4,
		TotalLen:       10,
	}

	params := &domain.PasswordGeneratorParams{
		Buffer:  buffer,
		Charset: charset,
		Chunk:   chunk,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		params.Iteration = uint64(i)
		generator.GeneratePassword(params)
	}
}
