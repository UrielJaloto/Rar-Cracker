package main

import (
	"fmt"
	"os"

	"github.com/UrielJaloto/Rar-Cracker/internal/config"
)

func main() {
	appConfig := config.New()

	configWarnings, err := appConfig.Validate()

	if len(configWarnings) > 0 {
		fmt.Println("(WARNINGS):")

		for _, warning := range configWarnings {
			fmt.Printf("  %s\n", warning)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s ", err)
		os.Exit(1)
	}

	appConfig.PrintFields()
}
