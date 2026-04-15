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
	fileHeaderType       = 0x02
	serviceHeaderType    = 0x03
	encryptionHeaderType = 0x04
	endArchiveHeaderType = 0x05
)

func ExtractEncryptionMetadata(reader io.ReadSeeker) (*domain.EncryptionMetadata, error) {
	if err := validateSignature(reader); err != nil {
		return nil, err
	}

	for {
		blockHeader, err := readBlockHeader(reader)
		if err != nil {
			return nil, err
		}

		if blockHeader.HeaderType == endArchiveHeaderType {
			return nil, errors.New("encryption header not found")
		}

		if (blockHeader.HeaderType == fileHeaderType || blockHeader.HeaderType == serviceHeaderType) && (blockHeader.HasExtraArea) {
			reader.Seek(blockHeader.BytesToReachExtraArea, io.SeekCurrent)

			var bytesProcessed int64
			for bytesProcessed < blockHeader.ExtraAreaSize {
				var extraAreaRecord *domain.ExtraAreaRecord

				if extraAreaRecord, err = readExtraArea(reader); err != nil {
					return nil, err
				}

				if extraAreaRecord.Type == 0x01 {
					return parseEncryptionHeader(reader, true)
				}

				reader.Seek(extraAreaRecord.BytesToEnd, io.SeekCurrent)
				bytesProcessed += extraAreaRecord.TotalSize
			}
			reader.Seek(-int64(blockHeader.BytesToReachExtraArea+blockHeader.ExtraAreaSize), io.SeekCurrent)
		}

		if blockHeader.HeaderType == encryptionHeaderType {
			return parseEncryptionHeader(reader, false)
		}

		if _, err := reader.Seek(blockHeader.BytesToReachNextBlock, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to skip header body: %w", err)
		}
	}
}

func readVarInt(reader io.Reader) (decodedValue uint64, bytesRead int64, err error) {
	var shift uint
	byteBuffer := make([]byte, 1)

	for {
		if _, err := io.ReadFull(reader, byteBuffer); err != nil {
			return 0, bytesRead, err
		}
		byteValue := byteBuffer[0]
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

func readBlockHeader(reader io.Reader) (blockHeader *domain.BlockHeader, err error) {
	blockHeader = &domain.BlockHeader{}

	headerCrc := make([]byte, 4)
	if _, err := io.ReadFull(reader, headerCrc); err != nil {
		return blockHeader, fmt.Errorf("failed to read CRC: %w", err)
	}

	headerSize, _, err := readVarInt(reader)
	if err != nil {
		return blockHeader, fmt.Errorf("failed to read header size: %w", err)
	}

	headerType, headerTypeLength, err := readVarInt(reader)
	if err != nil {
		return blockHeader, fmt.Errorf("failed to read header type: %w", err)
	}
	blockHeader.HeaderType = headerType
	headerBytesRead := headerTypeLength

	headerFlags, flagsBytesSize, err := readVarInt(reader)
	if err != nil {
		return blockHeader, fmt.Errorf("inconsistent header flags: %w", err)
	}
	headerBytesRead += flagsBytesSize

	var extraAreaSize uint64
	if blockHeader.HasExtraArea = (headerFlags & 0x0001) != 0; blockHeader.HasExtraArea {
		var extraSizeLength int64
		extraAreaSize, extraSizeLength, err = readVarInt(reader)
		if err != nil {
			return blockHeader, fmt.Errorf("inconsistent extra area size: %w", err)
		}
		headerBytesRead += extraSizeLength
	}
	blockHeader.ExtraAreaSize = int64(extraAreaSize)

	var dataAreaSize uint64
	if blockHeader.HasDataArea = (headerFlags & 0x0002) != 0; blockHeader.HasDataArea {
		var dataAreaSizeLength int64
		dataAreaSize, dataAreaSizeLength, err = readVarInt(reader)
		if err != nil {
			return blockHeader, fmt.Errorf("inconsistent data area size: %w", err)
		}
		headerBytesRead += dataAreaSizeLength
	}

	unprocessedHeaderBytes := int64(headerSize) - headerBytesRead
	if unprocessedHeaderBytes < 0 {
		return blockHeader, errors.New("inconsistent header size computation: negative remaining bytes")
	}
	blockHeader.BytesToReachExtraArea = unprocessedHeaderBytes - int64(extraAreaSize)
	blockHeader.BytesToReachNextBlock = unprocessedHeaderBytes + int64(dataAreaSize)

	return blockHeader, nil
}

func readExtraArea(reader io.Reader) (extraArea *domain.ExtraAreaRecord, err error) {
	extraArea = &domain.ExtraAreaRecord{}

	extraAreaSize, extraAreaSizeLength, err := readVarInt(reader)
	if err != nil {
		return nil, err
	}
	extraArea.TotalSize = int64(extraAreaSize) + extraAreaSizeLength

	extraAreaType, extraAreaTypeLength, err := readVarInt(reader)
	if err != nil {
		return nil, err
	}
	extraArea.Type = extraAreaType
	extraArea.BytesToEnd = int64(extraAreaSize) - extraAreaTypeLength

	return extraArea, nil
}

func parseEncryptionHeader(reader io.Reader, hasIV bool) (*domain.EncryptionMetadata, error) {
	if _, _, err := readVarInt(reader); err != nil {
		return nil, fmt.Errorf("failed to read encryption version header: %w", err)
	}

	encryptionFlags, _, err := readVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read encryption flags: %w", err)
	}
	usePasswordCheck := (encryptionFlags & 0x0001) != 0

	kdfCount := make([]byte, 1)
	_, err = io.ReadFull(reader, kdfCount)
	if err != nil {
		return nil, fmt.Errorf("failed to read KDF count: %w", err)
	}
	iterations := int(kdfCount[0])

	salt := make([]byte, 16)
	_, err = io.ReadFull(reader, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to read salt: %w", err)
	}

	if hasIV {
		iv := make([]byte, 16)
		if _, err := io.ReadFull(reader, iv); err != nil {
			return nil, err
		}

	}

	var passwordCheck []byte
	if usePasswordCheck {
		passwordCheck = make([]byte, 8)
		_, err = io.ReadFull(reader, passwordCheck)
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
