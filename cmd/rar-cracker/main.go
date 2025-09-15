package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/UrielJaloto/Rar-Cracker/internal/config"
)

func main() {
	appConfig := config.New()

	configWarnings, configErrors := appConfig.Setup()
	if len(configWarnings) > 0 {
		fmt.Println("(WARNINGS):")
		for _, warning := range configWarnings {
			fmt.Printf("    %s\n\n", warning)
		}
	}
	if configErrors != nil {
		fmt.Fprintln(os.Stderr, "(ERRORS):")
		errorLines := strings.Split(configErrors.Error(), "\n")
		for _, line := range errorLines {
			if line != "" {
				fmt.Fprintf(os.Stderr, "    %s\n\n", line)
			}
		}
		os.Exit(1)
	}

	appConfig.PrintFields()
}
