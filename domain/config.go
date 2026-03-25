package domain

import (
	"math/big"
	"unicode/utf8"
)

type Config struct {
	CharsetPath   string
	FilePath      string
	StateFilePath string
	KnownPart     string
	MaxLength     int
	Workers       int
	Charset       []rune
}

func (c *Config) CalculateTotalCombinations() *big.Int {
	charsetSize := big.NewInt(int64(len(c.Charset)))
	knownLen := utf8.RuneCountInString(c.KnownPart)
	remainingLength := int64(c.MaxLength - knownLen)

	if charsetSize.Sign() <= 0 || remainingLength <= 0 {
		return big.NewInt(0)
	}

	total := big.NewInt(0)
	exponent := big.NewInt(remainingLength)
	total.Exp(charsetSize, exponent, nil)
	return total
}
