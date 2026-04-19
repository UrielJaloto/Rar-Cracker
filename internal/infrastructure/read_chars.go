package infrastructure

import (
	"os"
	"strings"
)

func readChars(filePath string) ([]rune, error) {
	bytesContent, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	cleanString := strings.TrimSpace(string(bytesContent))
	return []rune(cleanString), nil
}

func readUniqueChars(filePath string) ([]rune, error) {
	allChars, err := readChars(filePath)
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
