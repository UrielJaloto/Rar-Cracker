package main

import (
	"fmt"
	"os"

	"github.com/UrielJaloto/Rar-Cracker/internal/config"
)

func main() {
	appConfig := config.New()

	configWarnings, configErrors := appConfig.Setup()
	if len(configWarnings) > 0 {
		fmt.Println(configWarnings)
	}
	if configErrors != nil {
		fmt.Fprintf(os.Stderr, "%s ", configErrors)
		os.Exit(1)
	}

	appConfig.PrintFields()
}
