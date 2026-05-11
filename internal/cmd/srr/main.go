package main

import (
	"github.com/UrielJaloto/surgical-rar-recovery/internal/infrastructure"
	"github.com/UrielJaloto/surgical-rar-recovery/internal/services"
)

func main() {
	configLoader := infrastructure.NewFlagLoader()
	configValidator := infrastructure.NewConfigValidator()
	ui := infrastructure.NewCli()
	parser := infrastructure.NewParser()

	application := services.NewApplication(configLoader, configValidator, ui, parser)

	application.Run()
}
