package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/UrielJaloto/Rar-Cracker/internal/domain"
	"github.com/UrielJaloto/Rar-Cracker/internal/infraestructure"
)

func NewCli() infraestructure.UiInterface {
	return &cli{}
}

type cli struct{}

func (c cli) ShowWarnings(warnings []string) {
	fmt.Println("(WARNINGS):")
	for _, w := range warnings {
		fmt.Printf("    %s\n", w)
	}
	fmt.Println()
}

func (c cli) ShowErrors(err error) {
	fmt.Fprintln(os.Stderr, "(ERRORS):")
	for line := range strings.SplitSeq(err.Error(), "\n") {
		if line != "" {
			fmt.Fprintf(os.Stderr, "    %s\n", line)
		}
	}
}

func (c cli) ShowConfiguration(config *domain.Config) {
	fmt.Printf("Loaded Configuration:\n")
	fmt.Printf("    Charset Path   : %s\n", config.CharsetPath)
	fmt.Printf("    File Path      : %s\n", config.FilePath)
	fmt.Printf("    Known Part     : %s\n", config.KnownPart)
	fmt.Printf("    State File     : %s\n", config.StateFilePath)
	fmt.Printf("    Max Length     : %d\n", config.MaxLength)
	fmt.Printf("    Workers        : %d\n", config.Workers)
}

func (c cli) ShowEncryptionMetadata(metaData *domain.EncryptionMetadata) {
	fmt.Printf("Encryption Metadata:\n")
	fmt.Printf("    Password Check     : %x\n", metaData.PasswordCheck)
	fmt.Printf("    Salt               : %x\n", metaData.Salt)
	fmt.Printf("    Iterations         : %d\n", metaData.Iterations)
	fmt.Printf("    UsePassword Check  : %t\n", metaData.UsePasswordCheck)
}
