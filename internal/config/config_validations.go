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
	WARNING_TRESHOLD = 250000000
	ERROR_TRESHOLD   = 1000000000000
)

func errorBuilder(errorMessages []string) error {
	if len(errorMessages) == 0 {
		return nil
	}

	var errs []error
	for _, message := range errorMessages {
		errs = append(errs, errors.New(message))
	}
	return errors.Join(errs...)
}

func (c *configFields) Setup() (warnings []string, err error) {
	var allWarnings, allErrors []string

	allErrors = append(allErrors, c.validateRequired()...)
	if len(allErrors) > 0 {
		return allWarnings, errorBuilder(allErrors)
	}

	allErrors = append(allErrors, c.validateFilePaths()...)
	if len(allErrors) > 0 {
		return allWarnings, errorBuilder(allErrors)
	}

	workersWarning, workersErrors := c.validateWorkers()
	allWarnings = append(allWarnings, workersWarning...)
	allErrors = append(allErrors, workersErrors...)

	if fileError := c.LoadCharset(); fileError != nil {
		allErrors = append(allErrors, fileError.Error())
		return allWarnings, errorBuilder(allErrors)
	}

	combinationsWarnings, combinationsErrors := c.validateCombinations()
	allWarnings = append(allWarnings, combinationsWarnings...)
	allErrors = append(allErrors, combinationsErrors...)
	if len(allErrors) > 0 {
		return allWarnings, errorBuilder(allErrors)
	}
	return allWarnings, nil
}

func (c *configFields) validateRequired() (errorsList []string) {
	if strings.TrimSpace(c.CharsetPath) == "" {
		errorsList = append(errorsList, "Charset is required")
	}
	if strings.TrimSpace(c.FilePath) == "" {
		errorsList = append(errorsList, "File is required")
	}

	return errorsList
}

func (c *configFields) validateFilePaths() (errorsList []string) {
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
			errorsList = append(errorsList, fmt.Sprintf("%s file not found: (%s)", name, path))
		}
	}
	return errorsList
}

func (c *configFields) validateWorkers() (warnings []string, errorsList []string) {
	if c.Workers <= 0 {
		errorsList = append(errorsList, fmt.Sprintf("Workers (%d) must be positive", c.Workers))
		return warnings, errorsList
	}

	availableCPU := runtime.NumCPU()
	allowed := int(float64(availableCPU) * 1.2)

	if c.Workers > allowed {
		errorsList = append(errorsList, fmt.Sprintf("Workers (%d) exceed number of CPUs (%d)", c.Workers, availableCPU))

	} else if c.Workers > availableCPU {
		warnings = append(warnings, fmt.Sprintf("Workers (%d) exceed 20%% above available CPUs (%d). Performance may degrade", c.Workers, availableCPU))
	}

	return warnings, errorsList
}

func (c *configFields) validateCombinations() (warnings []string, errorsList []string) {
	charsetLength := int64(len(c.Charset))
	if charsetLength <= 0 {
		errorsList = append(errorsList, ("Charset is empty, cannot calculate the number of combinations"))
		return warnings, errorsList
	}

	remainingLength := int64(c.MaxLenght - len(c.KnownPart))
	if remainingLength <= 0 {
		errorsList = append(errorsList, fmt.Sprintf("MaxLength (%d) cannot be smaller or equal then the KnownPart (%d)", c.MaxLenght, len(c.KnownPart)))
		return warnings, errorsList
	}

	totalCombinations := big.NewInt(0)
	charsetSize := big.NewInt(charsetLength)
	expoent := big.NewInt(remainingLength)
	errorTreashold := big.NewInt(ERROR_TRESHOLD)
	warningTreashold := big.NewInt(WARNING_TRESHOLD)

	totalCombinations.Exp(charsetSize, expoent, nil)

	if totalCombinations.Cmp(errorTreashold) >= 0 {
		errorsList = append(errorsList, fmt.Sprintf("Total combinations (%s) exceded ERROR_TREASHOLD (%d)", totalCombinations.String(), ERROR_TRESHOLD))

	} else if totalCombinations.Cmp(warningTreashold) >= 0 {
		warnings = append(warnings, fmt.Sprintf("Total combinations (%s) exceded WARNING_TREASHOLD (%d)", totalCombinations.String(), WARNING_TRESHOLD))
	}

	return warnings, errorsList
}
