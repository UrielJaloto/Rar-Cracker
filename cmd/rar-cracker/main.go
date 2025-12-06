package main

import (
	"fmt"
	"os"

	"github.com/UrielJaloto/Rar-Cracker/internal/config"
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
		fmt.Println("WARNINGS:")
		for _, w := range warnings {
			fmt.Printf("- %s\n", w)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERRORS:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting cracker for file: %s with %d workers\n", cfg.FilePath, cfg.Workers)
}
