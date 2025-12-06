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
	loader := config.NewLoader()
	cfg, err := loader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	validator := config.NewValidator()
	warnings, err := validator.Validate(cfg)

	if len(warnings) > 0 {
		fmt.Println("(WARNINGS):")
		for _, w := range warnings {
			fmt.Printf("    %s\n", w)
		}
		fmt.Println()
	}

	if err != nil {
		printErrors(err)
		os.Exit(1)
	}

	printConfiguration(cfg)

	file, err := os.Open(cfg.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	metadata, err := rar.ExtractMetadata(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing RAR metadata: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Metadata extracted successfully! (Encryption found: %v)\n", metadata != nil)
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
