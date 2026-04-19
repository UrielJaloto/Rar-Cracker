package main

import (
	"fmt"
	"os"

	"github.com/UrielJaloto/Rar-Cracker/internal/infrastructure"
)

func main() {
	configLoader := infrastructure.NewFlagLoader()
	validator := infrastructure.NewConfigValidator()
	userInterface := infrastructure.NewCli()
	extractor := infrastructure.NewParser()

	settings, err := configLoader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	warnings, err := validator.Validate(settings)
	if len(warnings) > 0 {
		userInterface.ShowWarnings(warnings)
	}
	if err != nil {
		userInterface.ShowErrors(err)
		os.Exit(1)
	}

	userInterface.ShowConfiguration(settings)

	file, err := os.Open(settings.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	encryptionMetadata, err := extractor.Extract(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing encryption metadata: %v\n", err)
		os.Exit(1)
	}
	println()

	userInterface.ShowEncryptionMetadata(encryptionMetadata)
}
