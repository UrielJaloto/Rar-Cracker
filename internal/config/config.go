package config

import (
	"flag"
	"fmt"

	"github.com/UrielJaloto/Rar-Cracker/internal/utils"
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
	RequiredParameters
	OptionalParameters
	Charset []rune
}

func New() *configFields {
	config := configFields{
		Charset: []rune{},
	}

	flag.StringVar(&config.CharsetPath, "charset", "", "Path for the charset (Required)")
	flag.StringVar(&config.FilePath, "file", "", "Path for the file (Required)")
	flag.StringVar(&config.KnownPart, "knownPart", "", "Known part of the password (Optional)")
	flag.StringVar(&config.StateFilePath, "stateFile", "./state-file.json", "Path for the state file (Optional)")
	flag.IntVar(&config.MaxLenght, "maxLength", 13, "Maximum length of the password (Optional)")
	flag.IntVar(&config.Workers, "workers", 1, "Number of workers")
	flag.Parse()

	return &config
}

func (c *configFields) LoadCharset() (err error) {
	c.Charset, err = utils.ReadChars(c.CharsetPath)
	return err
}

func (c *configFields) PrintFields() {
	fmt.Printf("Configurações carregadas:\n")
	fmt.Printf("    Charset Path   : %s\n", c.CharsetPath)
	fmt.Printf("    File Path      : %s\n", c.FilePath)
	fmt.Printf("    Known Part     : %s\n", c.KnownPart)
	fmt.Printf("    State File     : %s\n", c.StateFilePath)
	fmt.Printf("    Max Length     : %d\n", c.MaxLenght)
	fmt.Printf("    Workers        : %d\n", c.Workers)
	fmt.Printf("    Charset        : %q\n", c.Charset)
}
