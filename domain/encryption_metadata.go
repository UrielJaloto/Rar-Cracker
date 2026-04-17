package domain

import "io"

type MetadataExtractor interface {
	Extract(reader io.ReadSeeker) (*EncryptionMetadata, error)
}

type EncryptionMetadata struct {
	UsePasswordCheck bool
	Iterations       int
	Salt             []byte
	PasswordCheck    []byte
}
