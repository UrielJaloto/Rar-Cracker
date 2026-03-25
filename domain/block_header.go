package domain

type BlockHeader struct {
	HeaderType           uint64
	BytesToNextBlock     int64
	RemainingHeaderBytes int64
	ExtraAreaSize        int64
}
