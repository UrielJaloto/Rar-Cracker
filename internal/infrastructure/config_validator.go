package infrastructure

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/UrielJaloto/surgical-rar-recovery/internal/domain"
)

const (
	warningThreshold        = 250000000
	errorThreshold          = 1000000000
	workerUtilizationFactor = 1.2
)

type Validator struct{}

func NewConfigValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(settings *domain.Config) (validationReport domain.ValidationReport) {
	var errs []error

	if strings.TrimSpace(settings.CharsetPath) == "" {
		errs = append(errs, errors.New("charset path is required"))
	}

	if strings.TrimSpace(settings.FilePath) == "" {
		errs = append(errs, errors.New("file path is required"))
	}

	if len(settings.Charset) == 0 && settings.CharsetPath != "" {
		errs = append(errs, errors.New("charset is empty"))
	}

	if len(errs) > 0 {
		validationReport.Err = errors.Join(errs...)
		return validationReport
	}

	resourcesWarnings, err := v.validateResources(settings)
	validationReport.Warnings = append(validationReport.Warnings, resourcesWarnings...)
	if err != nil {
		errs = append(errs, err)
	}

	complexityWarning, complexityError := v.validateComplexity(settings)
	validationReport.Warnings = append(validationReport.Warnings, complexityWarning...)
	if complexityError != nil {
		errs = append(errs, complexityError)
	}

	validationReport.Err = errors.Join(errs...)
	return validationReport
}

func (v *Validator) validateResources(settings *domain.Config) (warnings []string, err error) {
	warnings = make([]string, 0)
	if settings.Workers <= 0 {
		err = fmt.Errorf("workers must be positive (got %d)", settings.Workers)
		return warnings, err
	}

	cpus := runtime.NumCPU()
	limit := int(float64(cpus) * workerUtilizationFactor)

	if settings.Workers > limit {
		err = fmt.Errorf("workers (%d) exceed CPU limit significantly (max allowed: %d)", settings.Workers, limit)
		return warnings, err
	} else if settings.Workers > cpus {
		warnings = append(warnings, fmt.Sprintf("workers (%d) exceed physical CPU count (%d). Performance may degrade.", settings.Workers, cpus))
	}

	return warnings, err
}

func (v *Validator) validateComplexity(settings *domain.Config) (warnings []string, err error) {
	knownLen := utf8.RuneCountInString(settings.KnownPart)
	remainingLength := settings.MaxLength - knownLen
	if remainingLength <= 0 {
		err = fmt.Errorf("max length (%d) cannot be smaller or equal to known part length (%d)", settings.MaxLength, len(settings.KnownPart))
		return warnings, err
	}

	total := settings.CalculateTotalCombinations()

	limitError := big.NewInt(errorThreshold)
	limitWarning := big.NewInt(warningThreshold)

	if total.Cmp(limitError) >= 0 {
		err = fmt.Errorf("total combinations (%s) exceeded the error threshold (%d)", total.String(), errorThreshold)
		return warnings, err
	}
	if total.Cmp(limitWarning) >= 0 {
		warnings = append(warnings, fmt.Sprintf("total combinations (%s) exceeded the warning threshold (%d)", total.String(), warningThreshold))
	}
	return warnings, err
}
