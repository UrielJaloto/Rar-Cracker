package config

import (
	"flag"
	"fmt"

	"github.com/UrielJaloto/Rar-Cracker/internal/domain"
	"github.com/UrielJaloto/Rar-Cracker/internal/infraestructure"
	"github.com/UrielJaloto/Rar-Cracker/internal/services/utils"
)

func NewLoader() infraestructure.ConfigLoaderInterface {
	return &FlagLoader{}
}

type FlagLoader struct{}

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

	var charset []rune
	var err error

	if charsetPath != "" {
		charset, err = utils.ReadUniqueChars(charsetPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read charset file: %w", err)
		}
	}

	return &domain.Config{
		CharsetPath:   charsetPath,
		FilePath:      filePath,
		StateFilePath: stateFile,
		KnownPart:     knownPart,
		MaxLength:     maxLength,
		Workers:       workers,
		Charset:       charset, // Agora ele vai preenchido para o Validator!
	}, nil
}
