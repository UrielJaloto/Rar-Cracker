package domain

type MetadataExtractor interface {
	Extract() (*EncryptionMetadata, error)
}

type EncryptionMetadata struct {
	UsePasswordCheck bool
	Iterations       int
	Salt             []byte
	PasswordCheck    []byte
}
