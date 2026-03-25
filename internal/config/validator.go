package config

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strings"

	"github.com/UrielJaloto/Rar-Cracker/domain"
)

const (
	warningThreshold        = 250000000
	errorThreshold          = 1000000000
	workerUtilizationFactor = 1.2
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(settings *domain.Config) ([]string, error) {
	var errs []error
	var warnings []string

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
		return nil, errors.Join(errs...)
	}

	if err := v.validatePaths(settings); err != nil {
		errs = append(errs, err)
	}

	warning, err := v.validateResources(settings)
	warnings = append(warnings, warning...)
	if err != nil {
		errs = append(errs, err)
	}

	complexityWarning, complexityError := v.validateComplexity(settings)
	warnings = append(warnings, complexityWarning...)
	if complexityError != nil {
		errs = append(errs, complexityError)
	}

	return warnings, errors.Join(errs...)
}

func (v *Validator) validatePaths(settings *domain.Config) error {
	paths := map[string]string{
		"Charset": settings.CharsetPath,
		"File":    settings.FilePath,
	}

	var errs []error
	for name, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("%s not found: %s", name, path))
		}
	}

	return errors.Join(errs...)
}

func (v *Validator) validateResources(settings *domain.Config) ([]string, error) {
	var warnings []string
	if settings.Workers <= 0 {
		return nil, fmt.Errorf("workers must be positive (got %d)", settings.Workers)
	}

	cpus := runtime.NumCPU()
	limit := int(float64(cpus) * workerUtilizationFactor)

	if settings.Workers > limit {
		return nil, fmt.Errorf("workers (%d) exceed CPU limit significantly (max allowed: %d)", settings.Workers, limit)
	} else if settings.Workers > cpus {
		warnings = append(warnings, fmt.Sprintf("workers (%d) exceed physical CPU count (%d). Performance may degrade", settings.Workers, cpus))
	}
	return warnings, nil
}

func (v *Validator) validateComplexity(settings *domain.Config) ([]string, error) {
	var warnings []string

	remainingLength := settings.MaxLength - len(settings.KnownPart)
	if remainingLength <= 0 {
		return nil, fmt.Errorf("max length (%d) cannot be smaller or equal to known part length (%d)", settings.MaxLength, len(settings.KnownPart))
	}

	total := settings.CalculateTotalCombinations()

	limitError := big.NewInt(errorThreshold)
	limitWarning := big.NewInt(warningThreshold)

	if total.Cmp(limitError) >= 0 {
		return nil, fmt.Errorf("total combinations (%s) exceeded the error threshold (%d)", total.String(), errorThreshold)
	}
	if total.Cmp(limitWarning) >= 0 {
		warnings = append(warnings, fmt.Sprintf("total combinations (%s) exceeded the warning threshold (%d)", total.String(), warningThreshold))
	}
	return warnings, nil
}
