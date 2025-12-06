package rar

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

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

func ExtractMetadata(reader io.ReadSeeker) (*domain.RarMetadata, error) {
	bufioReader := bufio.NewReader(reader)

	if err := validateSignature(bufioReader); err != nil {
		return nil, err
	}

	for {
		headerType, bodySize, err := readBlockHeader(bufioReader)
		if err != nil {
			return nil, err
		}

		if headerType == headerTypeEncryption {
			return &domain.RarMetadata{}, nil
		}

		if headerType == headerTypeEndArchive {
			return nil, errors.New("encryption header not found")
		}

		if _, err := reader.Seek(int64(bodySize), io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to skip header body: %w", err)
		}

		bufioReader.Reset(reader)
	}
}

func readBlockHeader(reader *bufio.Reader) (headerType uint64, bodySize uint64, err error) {
	headerCrc := make([]byte, 4)
	if _, err := io.ReadFull(reader, headerCrc); err != nil {
		return 0, 0, fmt.Errorf("failed to read CRC: %w", err)
	}

	rawHeaderSize, _, err := readVarInt(reader)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read header size: %w", err)
	}

	headerType, headerTypeSize, err := readVarInt(reader)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read header type: %w", err)
	}

	bytesToSkip := int64(rawHeaderSize) - int64(headerTypeSize)
	if bytesToSkip < 0 {
		return 0, 0, errors.New("inconsistent header size")
	}

	return headerType, uint64(bytesToSkip), nil
}

func validateSignature(reader io.Reader) error {
	signatureLen := len(rar5Signature)
	fileSignature := make([]byte, signatureLen)

	if _, err := io.ReadFull(reader, fileSignature); err != nil {
		return fmt.Errorf("failed to read signature: %w", err)
	}

	if !bytes.Equal(fileSignature, rar5Signature) {
		return errors.New("invalid signature: not a RAR5 file")
	}

	return nil
}

func readVarInt(reader io.ByteReader) (uint64, int, error) {
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
			return 0, bytesRead, errors.New("variable integer exceeded 64 bits")
		}
	}
}
