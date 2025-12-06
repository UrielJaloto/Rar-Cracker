package domain

type RarMetadata struct {
	Salt             []byte
	PasswordCheck    []byte
	UsePasswordCheck bool
	Iterations       int
}
