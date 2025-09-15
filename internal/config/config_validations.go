package config

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strings"
)

const (
	WARNING_TREASHOLD = 250000000
	ERROR_TREASHOLD   = 1000000000000
)

func (c *configFields) Validate() (warnings []string, err error) {
	var (
		allWarnings []string
		allErrors   []string
	)

	allErrors = append(allErrors, c.validateRequired()...)
	allErrors = append(allErrors, c.validateFilePaths()...)

	workersWarning, workersErrors := c.validateWorkers()
	allWarnings = append(allWarnings, workersWarning...)
	allErrors = append(allErrors, workersErrors...)

	combinationsWarnings, combinationsErrors := c.validateCombinations()
	allWarnings = append(allWarnings, combinationsWarnings...)
	allErrors = append(allErrors, combinationsErrors...)

	if len(allErrors) > 0 {
		errorMessage := "Errors found{\n  " + strings.Join(allErrors, "\n  ") + "\n}\nExecution blocked."
		return allWarnings, errors.New(errorMessage)
	}
	return allWarnings, nil
}

func (c *configFields) validateRequired() (errors []string) {
	if strings.TrimSpace(c.CharsetPath) == "" {
		errors = append(errors, "Charset is required")
	}
	if strings.TrimSpace(c.FilePath) == "" {
		errors = append(errors, "File is required")
	}

	return errors
}

func (c *configFields) validateFilePaths() (errors []string) {
	filePathsToValidate := map[string]string{
		"Charset": c.CharsetPath,
		"Archive": c.FilePath,
		"State":   c.StateFilePath,
	}

	for name, path := range filePathsToValidate {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("%s file not found: (%s)", name, path))
		}
	}
	return errors
}

func (c *configFields) validateWorkers() (warnings []string, errors []string) {
	if c.Workers <= 0 {
		errors = append(errors, fmt.Sprintf("Workers (%d) must be positive", c.Workers))
		return warnings, errors
	}

	availableCPU := runtime.NumCPU()
	allowed := int(float64(availableCPU) * 1.2)

	if c.Workers > allowed {
		errors = append(errors, fmt.Sprintf("Workers (%d) exceed 20%% above available CPUs (%d)", c.Workers, availableCPU))

	} else if c.Workers > availableCPU {
		warnings = append(warnings, fmt.Sprintf("Workers (%d) exceed number of CPUs (%d). Performance may degrade", c.Workers, availableCPU))
	}

	return warnings, errors
}

func (c *configFields) validateCombinations() (warnings []string, errors []string) {
	charsetLength := int64(len(c.Charset))
	if charsetLength <= 0 {
		errors = append(errors, ("Charset is empty, cannot calculate the number of combinations"))
		return warnings, errors
	}

	remainingLength := int64(c.MaxLenght - len(c.KnownPart))
	if remainingLength <= 0 {
		errors = append(errors, fmt.Sprintf("MaxLength (%d) cannot be smaller or equal then the KnownPart (%d)", c.MaxLenght, len(c.KnownPart)))
		return warnings, errors
	}

	totalCombinations := big.NewInt(0)
	charsetSize := big.NewInt(charsetLength)
	expoent := big.NewInt(remainingLength)
	errorTreashold := big.NewInt(ERROR_TREASHOLD)
	warningTreashold := big.NewInt(WARNING_TREASHOLD)

	totalCombinations.Exp(charsetSize, expoent, nil)

	if totalCombinations.Cmp(errorTreashold) >= 0 {
		errors = append(errors, fmt.Sprintf("Total combinations (%s) exceded ERROR_TREASHOLD (%d)", totalCombinations.String(), ERROR_TREASHOLD))

	} else if totalCombinations.Cmp(warningTreashold) >= 0 {
		warnings = append(warnings, fmt.Sprintf("Total combinations (%s) exceded WARNING_TREASHOLD (%d)", totalCombinations.String(), WARNING_TREASHOLD))
	}

	return warnings, errors
}
