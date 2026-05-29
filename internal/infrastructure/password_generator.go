package infrastructure

import "github.com/UrielJaloto/surgical-rar-recovery/internal/domain"

type PasswordGenerator struct{}

func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{}
}

func (pg *PasswordGenerator) GeneratePassword(params *domain.PasswordGeneratorParams) {
	dividend := params.Iteration
	divisor := uint64(len(params.Charset))
	knownLen := params.Chunk.TotalLen - params.Chunk.VariableLen

	for i := 0; i < params.Chunk.VariableLen; i++ {
		charsetIndex := dividend % divisor

		insertIndex := i
		if insertIndex >= params.Chunk.KnownPartIndex {
			insertIndex += knownLen
		}

		params.Buffer[insertIndex] = byte(params.Charset[charsetIndex])
		dividend = dividend / divisor
	}
}
