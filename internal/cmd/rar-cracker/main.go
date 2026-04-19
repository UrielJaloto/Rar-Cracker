package main

import (
	"fmt"
	"os"

	"github.com/UrielJaloto/Rar-Cracker/internal/services/config"
	"github.com/UrielJaloto/Rar-Cracker/internal/services/rar"
	"github.com/UrielJaloto/Rar-Cracker/internal/services/ui"
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

	userInterface := ui.NewCli()

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

	extractor := rar.NewParser()
	encryptionMetadata, err := extractor.Extract(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing encryption metadata: %v\n", err)
		os.Exit(1)
	}
	println()

	userInterface.ShowEncryptionMetadata(encryptionMetadata)
}
