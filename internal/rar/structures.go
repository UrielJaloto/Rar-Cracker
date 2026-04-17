package rar

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
