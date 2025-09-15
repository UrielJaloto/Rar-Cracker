package utils

import (
	"os"
)

func ReadChars(filePath string) ([]rune, error) {
	bytesContent, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	sliceChars := []rune(string(bytesContent))

	return sliceChars, nil
}
