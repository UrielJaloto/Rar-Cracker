package rar

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/UrielJaloto/Rar-Cracker/internal/domain"
)

type blockHeader struct {
	Type                  uint64
	HasExtraArea          bool
	HasDataArea           bool
	ExtraAreaSize         int64
	BytesToReachNextBlock int64
	BytesToReachExtraArea int64
}

type extraAreaRecord struct {
	TotalSize  int64
	Type       uint64
	BytesToEnd int64
}

var rar5Signature = []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x01, 0x00}

const (
	fileHeaderType       = 0x02
	serviceHeaderType    = 0x03
	encryptionHeaderType = 0x04
	endArchiveHeaderType = 0x05
)

type Parser struct {
	reader io.ReadSeeker
}

func NewParser(ioReader io.ReadSeeker) *Parser {
	return &Parser{
		reader: ioReader,
	}
}

func (p *Parser) Extract() (*domain.EncryptionMetadata, error) {
	if err := p.validateSignature(); err != nil {
		return nil, err
	}

	for {
		blockHeader, err := p.readBlockHeader()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, errors.New("nenhuma criptografia encontrada neste arquivo")
			}
			return nil, err
		}

		switch blockHeader.Type {
		case endArchiveHeaderType:
			return nil, errors.New("encryption header not found")

		case encryptionHeaderType:
			return p.parseEncryptionMetaData(false)

		case fileHeaderType, serviceHeaderType:
			if !blockHeader.HasExtraArea {
				break
			}

			p.reader.Seek(blockHeader.BytesToReachExtraArea, io.SeekCurrent)
			var bytesProcessed int64

			for bytesProcessed < blockHeader.ExtraAreaSize {
				var extraAreaRecord *extraAreaRecord

				if extraAreaRecord, err = p.readExtraArea(); err != nil {
					return nil, err
				}

				if extraAreaRecord.Type == 0x01 {
					return p.parseEncryptionMetaData(true)
				}

				p.reader.Seek(extraAreaRecord.BytesToEnd, io.SeekCurrent)
				bytesProcessed += extraAreaRecord.TotalSize
			}
			p.reader.Seek(-int64(blockHeader.BytesToReachExtraArea+blockHeader.ExtraAreaSize), io.SeekCurrent)
		}

		if _, err := p.reader.Seek(blockHeader.BytesToReachNextBlock, io.SeekCurrent); err != nil {
			return nil, fmt.Errorf("failed to skip header body: %w", err)
		}
	}
}

func (p *Parser) validateSignature() error {
	signatureLen := len(rar5Signature)
	fileSignature := make([]byte, signatureLen)

	if _, err := io.ReadFull(p.reader, fileSignature); err != nil {
		return fmt.Errorf("failed to read signature: %w", err)
	}

	if !bytes.Equal(fileSignature, rar5Signature) {
		return errors.New("invalid signature: not a RAR5 file")
	}

	return nil
}

func (p *Parser) readBlockHeader() (header *blockHeader, err error) {
	header = &blockHeader{}

	headerCrc := make([]byte, 4)
	if _, err := io.ReadFull(p.reader, headerCrc); err != nil {
		return header, fmt.Errorf("failed to read CRC: %w", err)
	}

	headerSize, _, err := readVarInt(p.reader)
	if err != nil {
		return header, fmt.Errorf("failed to read header size: %w", err)
	}

	headerType, headerTypeLength, err := readVarInt(p.reader)
	if err != nil {
		return header, fmt.Errorf("failed to read header type: %w", err)
	}
	header.Type = headerType
	headerBytesRead := headerTypeLength

	headerFlags, flagsBytesSize, err := readVarInt(p.reader)
	if err != nil {
		return header, fmt.Errorf("inconsistent header flags: %w", err)
	}
	headerBytesRead += flagsBytesSize

	var extraAreaSize uint64
	if header.HasExtraArea = (headerFlags & 0x0001) != 0; header.HasExtraArea {
		var extraSizeLength int64
		extraAreaSize, extraSizeLength, err = readVarInt(p.reader)
		if err != nil {
			return header, fmt.Errorf("inconsistent extra area size: %w", err)
		}
		headerBytesRead += extraSizeLength
	}
	header.ExtraAreaSize = int64(extraAreaSize)

	var dataAreaSize uint64
	if header.HasDataArea = (headerFlags & 0x0002) != 0; header.HasDataArea {
		var dataAreaSizeLength int64
		dataAreaSize, dataAreaSizeLength, err = readVarInt(p.reader)
		if err != nil {
			return header, fmt.Errorf("inconsistent data area size: %w", err)
		}
		headerBytesRead += dataAreaSizeLength
	}

	unprocessedHeaderBytes := int64(headerSize) - headerBytesRead
	if unprocessedHeaderBytes < 0 {
		return header, errors.New("inconsistent header size computation: negative remaining bytes")
	}
	bytesToReachExtraArea := unprocessedHeaderBytes - int64(extraAreaSize)
	if bytesToReachExtraArea < 0 {
		return header, errors.New("corrupted archive: extra area is larger then the header size")
	}

	header.BytesToReachExtraArea = bytesToReachExtraArea
	header.BytesToReachNextBlock = unprocessedHeaderBytes + int64(dataAreaSize)

	return header, nil
}

func (p *Parser) readExtraArea() (extraArea *extraAreaRecord, err error) {
	extraArea = &extraAreaRecord{}

	extraAreaSize, extraAreaSizeLength, err := readVarInt(p.reader)
	if err != nil {
		return nil, err
	}
	extraArea.TotalSize = int64(extraAreaSize) + extraAreaSizeLength

	extraAreaType, extraAreaTypeLength, err := readVarInt(p.reader)
	if err != nil {
		return nil, err
	}
	extraArea.Type = extraAreaType
	extraArea.BytesToEnd = int64(extraAreaSize) - extraAreaTypeLength

	return extraArea, nil
}

func (p *Parser) parseEncryptionMetaData(hasIV bool) (*domain.EncryptionMetadata, error) {
	if _, _, err := readVarInt(p.reader); err != nil {
		return nil, fmt.Errorf("failed to read encryption version header: %w", err)
	}

	encryptionFlags, _, err := readVarInt(p.reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read encryption flags: %w", err)
	}
	usePasswordCheck := (encryptionFlags & 0x0001) != 0

	kdfCount := make([]byte, 1)
	_, err = io.ReadFull(p.reader, kdfCount)
	if err != nil {
		return nil, fmt.Errorf("failed to read KDF count: %w", err)
	}
	iterations := 1 << kdfCount[0]

	salt := make([]byte, 16)
	_, err = io.ReadFull(p.reader, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to read salt: %w", err)
	}

	if hasIV {
		iv := make([]byte, 16)
		if _, err := io.ReadFull(p.reader, iv); err != nil {
			return nil, err
		}

	}

	var passwordCheck []byte
	if usePasswordCheck {
		passwordCheck = make([]byte, 12)
		_, err = io.ReadFull(p.reader, passwordCheck)
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

func readVarInt(reader io.Reader) (decodedValue uint64, bytesRead int64, err error) {
	var shift uint
	var byteBuffer [1]byte
	for {
		if _, err := io.ReadFull(reader, byteBuffer[:]); err != nil {
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
