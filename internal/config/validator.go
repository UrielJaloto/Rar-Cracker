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
	WarningThreshold = 250000000
	ErrorThreshold   = 1000000000000
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(cfg *domain.Config) ([]string, error) {
	var errs []error
	var warnings []string

	if strings.TrimSpace(cfg.CharsetPath) == "" {
		errs = append(errs, errors.New("charset path is required"))
	}
	if strings.TrimSpace(cfg.FilePath) == "" {
		errs = append(errs, errors.New("file path is required"))
	}
	if len(cfg.Charset) == 0 && cfg.CharsetPath != "" {
		errs = append(errs, errors.New("charset is empty"))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	if err := v.validatePaths(cfg); err != nil {
		errs = append(errs, err)
	}

	w, e := v.validateResources(cfg)
	warnings = append(warnings, w...)
	if e != nil {
		errs = append(errs, e)
	}

	wComb, eComb := v.validateComplexity(cfg)
	warnings = append(warnings, wComb...)
	if eComb != nil {
		errs = append(errs, eComb)
	}

	return warnings, errors.Join(errs...)
}

func (v *Validator) validatePaths(cfg *domain.Config) error {
	paths := map[string]string{
		"Charset": cfg.CharsetPath,
		"File":    cfg.FilePath,
	}

	var errs []error
	for name, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("%s not found: %s", name, path))
		}
	}

	return errors.Join(errs...)
}

func (v *Validator) validateResources(cfg *domain.Config) ([]string, error) {
	var warnings []string
	if cfg.Workers <= 0 {
		return nil, fmt.Errorf("workers must be positive (got %d)", cfg.Workers)
	}

	cpus := runtime.NumCPU()
	limit := int(float64(cpus) * 1.2)

	if cfg.Workers > limit {
		return nil, fmt.Errorf("workers (%d) exceed CPU limit significantly (max allowed: %d)", cfg.Workers, limit)
	} else if cfg.Workers > cpus {
		warnings = append(warnings, fmt.Sprintf("workers (%d) exceed physical CPU count (%d). Performance may degrade", cfg.Workers, cpus))
	}
	return warnings, nil
}

func (v *Validator) validateComplexity(cfg *domain.Config) ([]string, error) {
	var warnings []string

	remainingLength := cfg.MaxLength - len(cfg.KnownPart)
	if remainingLength <= 0 {
		return nil, fmt.Errorf("max length (%d) cannot be smaller or equal to known part length (%d)", cfg.MaxLength, len(cfg.KnownPart))
	}

	total := cfg.CalculateTotalCombinations()

	limitErr := big.NewInt(ErrorThreshold)
	limitWarn := big.NewInt(WarningThreshold)

	if total.Cmp(limitErr) >= 0 {
		return nil, fmt.Errorf("total combinations (%s) exceeded ERROR_THRESHOLD (%d)", total.String(), ErrorThreshold)
	}
	if total.Cmp(limitWarn) >= 0 {
		warnings = append(warnings, fmt.Sprintf("total combinations (%s) exceeded WARNING_THRESHOLD (%d)", total.String(), WarningThreshold))
	}
	return warnings, nil
}
