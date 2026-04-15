package domain

type BlockHeader struct {
	HeaderType            uint64
	BytesToReachNextBlock int64
	BytesToReachExtraArea int64
	ExtraAreaSize         int64
	HasExtraArea          bool
	HasDataArea           bool
}
