package config

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func (c *configFields) ValidateFields() error {
	var errorsList []string

	errorsList = append(errorsList, c.validateRequired()...)
	errorsList = append(errorsList, c.validateFilePaths()...)
	errorsList = append(errorsList, c.validateWorkers())

	if len(errorsList) > 0 {
		errorMessage := "Errors found: " + strings.Join(errorsList, ", ")

		return errors.New(errorMessage)
	}
	return nil
}

func (c configFields) validateRequired() []string {
	var errorsList []string

	if strings.TrimSpace(c.CharsetPath) == "" {

		errorsList = append(errorsList, "Charset is required")
	}

	if strings.TrimSpace(c.FilePath) == "" {
		errorsList = append(errorsList, "File is required")
	}
	return errorsList
}

func (c configFields) validateFilePaths() []string {
	var errorsList []string

	filePathsToValidate := []string{
		c.CharsetPath,
		c.FilePath,
		c.StateFilePath,
	}

	for _, filePath := range filePathsToValidate {
		if filePath == "" {
			continue

		}
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			errorsList = append(errorsList, fmt.Sprintf("File not found: %s", filePath))
		}
	}
	return errorsList
}

func (c configFields) validateWorkers() string {
	availableCPU := runtime.NumCPU()

	allowed := int(float64(availableCPU) * 1.2)

	if c.Workers > allowed {
		err := fmt.Sprintf("Workers (%d) exceed 20%% above available CPUs (%d).", c.Workers, availableCPU)
		return err
	}

	if c.Workers > availableCPU {
		fmt.Printf("Workers (%d) exceed number of CPUs (%d). Performance may degrade.\n", c.Workers, availableCPU)
	}

	return ""
}
