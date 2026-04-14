package domain

type BlockHeader struct {
	HeaderType            uint64
	BytesToReachNextBlock int64
	RemainingHeaderBytes  int64
}
