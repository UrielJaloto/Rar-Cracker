package main

import (
	"os"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/infrastructure"
	"github.com/UrielJaloto/surgical-rar-recovery/internal/services"
)

func main() {
	settingsBuilder := services.NewSettingsBuilder(
		infrastructure.NewFlagLoader(),
		infrastructure.NewConfigValidator(),
	)

	ui := infrastructure.NewCli()

	recoveryEngine := services.NewRecoveryEngine(
		infrastructure.NewParser(),
		infrastructure.NewPbkdf2KeyStretcher(),
	)

	application := services.NewApplication(
		settingsBuilder,
		ui,
		recoveryEngine,
	)

	err := application.Run()
	if err != nil {
		os.Exit(1)
	}
}
