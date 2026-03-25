package config

import (
	"flag"

	"github.com/UrielJaloto/Rar-Cracker/domain"
)

type FlagLoader struct{}

func NewLoader() *FlagLoader {
	return &FlagLoader{}
}

func (f *FlagLoader) Load() (*domain.Config, error) {
	var charsetPath, filePath, knownPart, stateFile string
	var maxLength, workers int

	flag.StringVar(&charsetPath, "charset", "", "Path for the charset (Required)")
	flag.StringVar(&filePath, "file", "", "Path for the file (Required)")
	flag.StringVar(&knownPart, "known-part", "", "Known part of the password")
	flag.StringVar(&stateFile, "state-file", "./state-file.json", "Path for state file")
	flag.IntVar(&maxLength, "max-length", 13, "Maximum password length")
	flag.IntVar(&workers, "workers", 1, "Number of workers")

	flag.Parse()

	return &domain.Config{
		CharsetPath:   charsetPath,
		FilePath:      filePath,
		StateFilePath: stateFile,
		KnownPart:     knownPart,
		MaxLength:     maxLength,
		Workers:       workers,
	}, nil
}
