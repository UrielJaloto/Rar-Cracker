package services

import (
	"os"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

type Application struct {
	settingsBuilder *SettingsBuilder
	ui              UiInterface
	recoveryEngine  *RecoveryEngine
}

func NewApplication(settingsBuilder *SettingsBuilder, ui UiInterface, recoveryEngine *RecoveryEngine) Application {
	return Application{settingsBuilder, ui, recoveryEngine}
}

func (app Application) Run() (err error) {
	settings, validationReport := app.settingsBuilder.Build()

	if len(validationReport.Warnings) > 0 {
		app.ui.ShowWarnings(validationReport.Warnings)
	}
	if validationReport.Err != nil {
		app.ui.ShowErrors(validationReport.Err)
		err = validationReport.Err
		return err
	}
	app.ui.ShowConfiguration(settings)

	var encryptionMetadata *domain.EncryptionMetadata
	encryptionMetadata, err = app.recoveryEngine.ParseMetadata(settings)
	if err != nil {
		app.ui.ShowErrors(err)
		return err
	}

	app.ui.ShowEncryptionMetadata(encryptionMetadata)
	return err
}

type SettingsBuilder struct {
	configLoader    ConfigLoaderInterface
	configValidator ConfigValidatorInterface
}

func NewSettingsBuilder(configLoader ConfigLoaderInterface, configValidator ConfigValidatorInterface) (settingsBuilder *SettingsBuilder) {
	settingsBuilder = &SettingsBuilder{configLoader, configValidator}
	return settingsBuilder
}

func (sb SettingsBuilder) Build() (settings *domain.Config, validationReport domain.ValidationReport) {
	settings, validationReport.Err = sb.configLoader.Load()
	if validationReport.Err != nil {
		return settings, validationReport
	}

	validationReport = sb.configValidator.Validate(settings)
	return settings, validationReport
}

type RecoveryEngine struct {
	parser       ParserInterface
	KeyStretcher KeyStretcherInterface
}

func NewRecoveryEngine(parser ParserInterface, KeyStretcher KeyStretcherInterface) (recoveryEngine *RecoveryEngine) {
	recoveryEngine = &RecoveryEngine{parser, KeyStretcher}
	return recoveryEngine
}

func (re RecoveryEngine) ParseMetadata(settings *domain.Config) (encryptionMetadata *domain.EncryptionMetadata, err error) {
	var file *os.File
	file, err = os.Open(settings.FilePath)
	if err != nil {
		return encryptionMetadata, err
	}
	defer file.Close()

	encryptionMetadata, err = re.parser.Extract(file)
	if err != nil {
		return encryptionMetadata, err
	}

	return encryptionMetadata, err
}

func (re RecoveryEngine) Recovery() (password string) {
	//TODO
	return
}
