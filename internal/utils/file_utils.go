package utils

import (
	"os"
	"strings"
)

func ReadChars(filePath string) ([]rune, error) {
	bytesContent, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	cleanString := strings.TrimSpace(string(bytesContent))
	return []rune(cleanString), nil
}

func ReadUniqueChars(filePath string) ([]rune, error) {
	allChars, err := ReadChars(filePath)
	if err != nil {
		return nil, err
	}

	seenChars := make(map[rune]bool)
	var uniqueChars []rune
	for _, char := range allChars {
		if !seenChars[char] {
			seenChars[char] = true
			uniqueChars = append(uniqueChars, char)
		}
	}

	return uniqueChars, nil
}
