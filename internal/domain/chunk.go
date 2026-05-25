package domain

type Chunk struct {
	VariableLen    int
	KnownPartIndex int
	TotalLen       int
	StartIndex     uint64
	EndIndex       uint64
}
