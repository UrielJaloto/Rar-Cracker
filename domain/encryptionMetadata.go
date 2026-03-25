package domain

type EncryptionMetadata struct {
	UsePasswordCheck bool
	Iterations       int
	Salt             []byte
	PasswordCheck    []byte
}
