package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/UrielJaloto/Rar-Cracker/domain"
	"github.com/UrielJaloto/Rar-Cracker/internal/config"
	"github.com/UrielJaloto/Rar-Cracker/internal/rar"
)

func main() {
	configLoader := config.NewLoader()
	settings, err := configLoader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	validator := config.NewValidator()
	warnings, err := validator.Validate(settings)

	if len(warnings) > 0 {
		printWarnings(warnings)
	}

	if err != nil {
		printErrors(err)
		os.Exit(1)
	}

	printConfiguration(settings)

	file, err := os.Open(settings.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var extractor domain.MetadataExtractor = rar.NewParser()
	encryptionMetadata, err := extractor.Extract(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing encryption metadata: %v\n", err)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing encryption metadata: %v\n", err)
		os.Exit(1)
	}
	println()

	printEncryptionMetadata(encryptionMetadata)
}

func printWarnings(warnings []string) {
	fmt.Println("(WARNINGS):")
	for _, w := range warnings {
		fmt.Printf("    %s\n", w)
	}
	fmt.Println()
}

func printErrors(err error) {
	fmt.Fprintln(os.Stderr, "(ERRORS):")
	for line := range strings.SplitSeq(err.Error(), "\n") {
		if line != "" {
			fmt.Fprintf(os.Stderr, "    %s\n", line)
		}
	}
}

func printConfiguration(c *domain.Config) {
	fmt.Printf("Loaded Configuration:\n")
	fmt.Printf("    Charset Path   : %s\n", c.CharsetPath)
	fmt.Printf("    File Path      : %s\n", c.FilePath)
	fmt.Printf("    Known Part     : %s\n", c.KnownPart)
	fmt.Printf("    State File     : %s\n", c.StateFilePath)
	fmt.Printf("    Max Length     : %d\n", c.MaxLength)
	fmt.Printf("    Workers        : %d\n", c.Workers)
}

func printEncryptionMetadata(e *domain.EncryptionMetadata) {
	fmt.Printf("Encryption Metadata:\n")
	fmt.Printf("    Password Check     : %x\n", e.PasswordCheck)
	fmt.Printf("    Salt               : %x\n", e.Salt)
	fmt.Printf("    Iterations         : %d\n", e.Iterations)
	fmt.Printf("    UsePassword Check  : %t\n", e.UsePasswordCheck)
}
