package rar

import (
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

func ExtractEncryptionMetadata(reader io.ReadSeeker) (*domain.EncryptionMetadata, error) {
	if err := validateSignature(reader); err != nil {
		return nil, err
	}

	for {
		headerType, bodySize, err := readBlockHeader(reader)
		if err != nil {
			return nil, err
		}

		if headerType == headerTypeEndArchive {
			return nil, errors.New("encryption header not found")
		}

		if headerType == headerTypeEncryption {
			return parseEncryptionHeader(reader)
		}

		if _, err := reader.Seek(int64(bodySize), io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to skip header body: %w", err)
		}
	}
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

func readVarInt(reader io.Reader) (decodedValue uint64, bytesRead int, err error) {
	var shift uint
	Buffer := make([]byte, 1)

	for {
		if _, err := reader.Read(Buffer); err != nil {
			return 0, bytesRead, err
		}
		byteValue := Buffer[0]
		bytesRead++

		// & 0x7F (01111111): Strips the 8th bit (continuation flag), keeping only the 7 data bits.
		// << shift: Moves these 7 bits to their correct position in the final number.
		// |= (OR): Safely merges the shifted bits into decodedValue without overwriting existing 1s.
		decodedValue |= uint64((byteValue & 0x7F)) << shift

		if (byteValue & 0x80) == 0 {
			return decodedValue, bytesRead, nil
		}

		shift += 7
		if shift > 64 {
			return 0, bytesRead, errors.New("variable integer exceeded 64 bits")
		}
	}
}

func readBlockHeader(reader io.Reader) (headerType uint64, bytesToNextBlock uint64, err error) {
	headerCrc := make([]byte, 4)
	if _, err := io.ReadFull(reader, headerCrc); err != nil {
		return 0, 0, fmt.Errorf("failed to read CRC: %w", err)
	}

	rawHeaderSize, _, err := readVarInt(reader)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read header size: %w", err)
	}

	headerType, typeBytes, err := readVarInt(reader)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read header type: %w", err)
	}

	headerFlags, flagsBytes, err := readVarInt(reader)
	if err != nil {
		return 0, 0, errors.New("inconsistent header Flags")
	}

	headerBytesRead := int64(typeBytes + flagsBytes)
	var dataAreaSize uint64

	if (headerFlags & 0x0001) != 0 {
		_, extraBytes, err := readVarInt(reader)
		if err != nil {
			return 0, 0, errors.New("inconsistent Extra Area Size")
		}
		headerBytesRead += int64(extraBytes)
	}

	if (headerFlags & 0x0002) != 0 {
		var dataBytes int
		dataAreaSize, dataBytes, err = readVarInt(reader)
		if err != nil {
			return 0, 0, errors.New("inconsistent Data Area Size")
		}
		headerBytesRead += int64(dataBytes)
	}

	remainingHeaderBytes := int64(rawHeaderSize) - headerBytesRead
	if remainingHeaderBytes < 0 {
		return 0, 0, errors.New("inconsistent header size computation")
	}

	bytesToNextBlock = uint64(remainingHeaderBytes) + dataAreaSize

	return headerType, bytesToNextBlock, nil
}

func parseEncryptionHeader(reader io.Reader) (*domain.EncryptionMetadata, error) {
	KDFCount := make([]byte, 1)
	salt := make([]byte, 16)
	passwordCheck := make([]byte, 8)

	if _, _, err := readVarInt(reader); err != nil {
		return nil, fmt.Errorf("failed to read encryption version header: %w", err)
	}

	encryptionFlags, _, err := readVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read encryption flags: %w", err)
	}
	usePasswordCheck := (encryptionFlags & 0x01) != 0

	_, err = io.ReadFull(reader, KDFCount)
	if err != nil {
		return nil, fmt.Errorf("failed to read KDF count: %w", err)
	}
	iterations := int(KDFCount[0])

	_, err = reader.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to read salt: %w", err)
	}

	if usePasswordCheck {
		_, err = reader.Read(passwordCheck)
		if err != nil {
			return nil, fmt.Errorf("failed to read password check: %w", err)
		}
	}

	return &domain.EncryptionMetadata{
		UsePasswordCheck: usePasswordCheck,
		Iterations:       iterations,
		Salt:             salt,
		PasswordCheck:    passwordCheck,
	}, nil
}
