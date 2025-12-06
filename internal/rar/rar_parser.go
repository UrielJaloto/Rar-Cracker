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

const (
	headerTypeMainArchive = 0x01
	headerTypeFile        = 0x02
	headerTypeService     = 0x03
	headerTypeEncryption  = 0x04
	headerTypeEndArchive  = 0x05
)

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

	for {
		headerCrc := make([]byte, 4)
		if _, err := io.ReadFull(fileReader, headerCrc); err != nil {
			return nil, fmt.Errorf("falha ao ler CRC do cabeçalho: %w", err)
		}

		headerSize, _, err := readVariableLengthInteger(fileReader)
		if err != nil {
			return nil, fmt.Errorf("falha ao ler tamanho do cabeçalho: %w", err)
		}

		headerType, bytesReadForHeaderType, err := readVariableLengthInteger(fileReader)
		if err != nil {
			return nil, fmt.Errorf("falha ao ler tipo do cabeçalho: %w", err)
		}

		if headerType == headerTypeEncryption {
			return &domain.RarMetadata{}, nil
		}

		if headerType == headerTypeEndArchive {
			return nil, errors.New("cabeçalho de criptografia não encontrado antes do fim do arquivo")
		}

		bytesToSkip := int64(headerSize) - int64(bytesReadForHeaderType)

		if bytesToSkip < 0 {
			return nil, errors.New("tamanho do cabeçalho inconsistente")
		}

		if _, err := file.Seek(bytesToSkip, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("falha ao pular conteúdo do cabeçalho: %w", err)
		}

		fileReader.Reset(file)
	}
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

func readVariableLengthInteger(reader *bufio.Reader) (uint64, int, error) {
	var value uint64
	var shift uint
	var bytesRead int

	for {
		byteValue, err := reader.ReadByte()
		if err != nil {
			return 0, bytesRead, err
		}
		bytesRead++

		value |= uint64(byteValue&0x7F) << shift

		if byteValue&0x80 == 0 {
			return value, bytesRead, nil
		}

		shift += 7
		if shift > 64 {
			return 0, bytesRead, errors.New("inteiro de tamanho variável excedeu limite de 64 bits")
		}
	}
}
