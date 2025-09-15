package config

import (
	"flag"
	"fmt"
)

type RequiredParameters struct {
	CharsetPath string
	FilePath    string
}

type OptionalParameters struct {
	StateFilePath string
	KnownPart     string
	MaxLenght     int
	Workers       int
}

type configFields struct {
	*RequiredParameters
	*OptionalParameters
	Charset []rune
}

func (c configFields) PrintFields() {
	fmt.Printf("Configurações carregadas:\n")
	fmt.Printf("  Charset Path   : %s\n", c.CharsetPath)
	fmt.Printf("  File Path      : %s\n", c.FilePath)
	fmt.Printf("  Known Part     : %s\n", c.KnownPart)
	fmt.Printf("  State File     : %s\n", c.StateFilePath)
	fmt.Printf("  Max Length     : %d\n", c.MaxLenght)
	fmt.Printf("  Workers        : %d\n", c.Workers)
}

func New() *configFields {
	config := configFields{
		RequiredParameters: &RequiredParameters{},
		OptionalParameters: &OptionalParameters{},
		Charset:            []rune{},
	}

	flag.StringVar(&config.CharsetPath, "Charset", "", "Path for the charset (Required)")
	flag.StringVar(&config.FilePath, "File", "", "Path for the file (Required)")
	flag.StringVar(&config.KnownPart, "KnownPart", "", "Known part of the password (Optional)")
	flag.StringVar(&config.StateFilePath, "StateFile", "./state-file.json", "Path for the state file (Optional)")
	flag.IntVar(&config.MaxLenght, "MaxLength", 13, "Maximum length of the password (Optional)")
	flag.IntVar(&config.Workers, "Workers", 1, "Number of workers")
	flag.Parse()

	return &config
}
