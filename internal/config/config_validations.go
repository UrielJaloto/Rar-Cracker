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

func warningBuilder(warnings []string) string {
	if len(warnings) == 0 {
		return ""
	}

	var stringBuilder strings.Builder
	stringBuilder.WriteString("(WARNINGS):\n    ")
	stringBuilder.WriteString(strings.Join(warnings, "\n    ") + "\n")

	return stringBuilder.String()
}

func errorBuilder(errorsList []string) error {
	if len(errorsList) == 0 {
		return nil
	}

	var stringBuilder strings.Builder
	stringBuilder.WriteString("(ERRORS):\n    ")
	stringBuilder.WriteString(strings.Join(errorsList, "\n    ") + "\n")
	stringBuilder.WriteString("\nExecution blocked.")

	return errors.New(stringBuilder.String())
}

func (c *configFields) Setup() (warnings string, err error) {
	var allWarnings, allErrors []string

	allErrors = append(allErrors, c.validateRequired()...)

	allErrors = append(allErrors, c.validateFilePaths()...)

	workersWarning, workersErrors := c.validateWorkers()
	allWarnings = append(allWarnings, workersWarning...)
	allErrors = append(allErrors, workersErrors...)

	if fileError := c.LoadCharset(); fileError != nil {
		allErrors = append(allErrors, fileError.Error())
		return warningBuilder(allWarnings), errorBuilder(allErrors)
	}

	combinationsWarnings, combinationsErrors := c.validateCombinations()
	allWarnings = append(allWarnings, combinationsWarnings...)
	allErrors = append(allErrors, combinationsErrors...)

	if len(allErrors) > 0 {

		return warningBuilder(allWarnings), errorBuilder(allErrors)
	}
	return warningBuilder(allWarnings), nil
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
	errorTreashold := big.NewInt(ERROR_TREASHOLD)
	warningTreashold := big.NewInt(WARNING_TREASHOLD)

	totalCombinations.Exp(charsetSize, expoent, nil)

	if totalCombinations.Cmp(errorTreashold) >= 0 {
		errorsList = append(errorsList, fmt.Sprintf("Total combinations (%s) exceded ERROR_TREASHOLD (%d)", totalCombinations.String(), ERROR_TREASHOLD))

	} else if totalCombinations.Cmp(warningTreashold) >= 0 {
		warnings = append(warnings, fmt.Sprintf("Total combinations (%s) exceded WARNING_TREASHOLD (%d)", totalCombinations.String(), WARNING_TREASHOLD))
	}

	return warnings, errorsList
}
