package config

import (
	"flag"

	"github.com/UrielJaloto/Rar-Cracker/domain"
	"github.com/UrielJaloto/Rar-Cracker/internal/utils" // Importa o utils
)

type FlagLoader struct{}

func NewLoader() *FlagLoader {
	return &FlagLoader{}
}

func (l *FlagLoader) Load() (*domain.Config, error) {
	var charsetPath, filePath, knownPart, stateFile string
	var maxLength, workers int

	flag.StringVar(&charsetPath, "charset", "", "Path for the charset (Required)")
	flag.StringVar(&filePath, "file", "", "Path for the file (Required)")
	flag.StringVar(&knownPart, "knownPart", "", "Known part of the password")
	flag.StringVar(&stateFile, "stateFile", "./state-file.json", "Path for state file")
	flag.IntVar(&maxLength, "maxLength", 13, "Maximum password length")
	flag.IntVar(&workers, "workers", 1, "Number of workers")

	flag.Parse()

	var charset []rune
	var err error

	if charsetPath != "" {
		charset, err = utils.ReadChars(charsetPath)
		if err != nil {
			return nil, err
		}
	}

	return &domain.Config{
		CharsetPath:   charsetPath,
		FilePath:      filePath,
		StateFilePath: stateFile,
		KnownPart:     knownPart,
		MaxLength:     maxLength,
		Workers:       workers,
		Charset:       charset,
	}, nil
}
