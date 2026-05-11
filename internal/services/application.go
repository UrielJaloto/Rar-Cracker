package services

import (
	"fmt"
	"os"
)

type Application struct {
	configLoader    ConfigLoaderInterface
	configValidator ConfigValidatorInterface
	ui              UiInterface
	parser          ParserInterface
}

func NewApplication(configLoader ConfigLoaderInterface, configValidator ConfigValidatorInterface, ui UiInterface, parser ParserInterface) Application {
	return Application{configLoader, configValidator, ui, parser}
}

func (app Application) Run() {
	settings, err := app.configLoader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	warnings, err := app.configValidator.Validate(settings)
	if len(warnings) > 0 {
		app.ui.ShowWarnings(warnings)
	}
	if err != nil {
		app.ui.ShowErrors(err)
		os.Exit(1)
	}

	app.ui.ShowConfiguration(settings)

	file, err := os.Open(settings.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	encryptionMetadata, err := app.parser.Extract(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing encryption metadata: %v\n", err)
		os.Exit(1)
	}
	println()

	app.ui.ShowEncryptionMetadata(encryptionMetadata)
}
