package utils

import (
	"os"
)

type FileReader struct{}

func NewFileReader() *FileReader {
	return &FileReader{}
}

func (f *FileReader) ReadChars(filePath string) ([]rune, error) {
	bytesContent, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	sliceChars := []rune(string(bytesContent))

	return sliceChars, nil
}
