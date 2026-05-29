package domain

type PasswordGeneratorParams struct {
	Buffer    []byte
	Charset   []rune
	Iteration uint64
	VarLen    int
}
