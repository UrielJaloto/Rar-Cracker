package domain

type RarMetadata struct {
	Salt        []byte
	PswCheck    []byte
	UsePswCheck bool
	Iterations  int
}
