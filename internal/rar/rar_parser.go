package rar

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/UrielJaloto/Rar-Cracker/domain"
)

var rar5Signature = []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x01, 0x00}

func ExtractMetadata(filePath string) (*domain.RarMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir arquivo: %w", err)
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)

	if err := checkRar5Signature(fileReader); err != nil {
		return nil, err
	}

	return &domain.RarMetadata{}, nil
}

func checkRar5Signature(reader io.Reader) error {
	signatureLen := len(rar5Signature)
	fileSignature := make([]byte, signatureLen)

	if _, err := io.ReadFull(reader, fileSignature); err != nil {
		return fmt.Errorf("não foi possível ler a assinatura: %w", err)
	}

	if !bytes.Equal(fileSignature, rar5Signature) {
		return errors.New("assinatura inválida: o arquivo não é um RAR5")
	}

	return nil
}
