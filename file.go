package extended_fs

import (
	"fmt"
	"strings"
)

// WriteFile will write a file content to a new file or override an existing one
func WriteFile(filePath string, fileContent string) error {
	for _, handler := range Handlers {
		if handler.CanHandle(filePath) {
			if !handler.AllowWrite() {
				return fmt.Errorf("writing is not allowed")
			}
			return handler.WriteFile(filePath, fileContent)
		}
	}
	return fmt.Errorf("unable to handle WriteFile")
}

// AppendToFile will append file content to an existing one
func AppendToFile(filePath string, fileContent string, onlyIfMissing bool) error {
	if onlyIfMissing {
		content, err := ReadFile(filePath)
		if err != nil {
			return err
		}
		if strings.Contains(content, fileContent) {
			return nil
		}
	}

	for _, handler := range Handlers {
		if handler.CanHandle(filePath) {
			if !handler.AllowWrite() {
				return fmt.Errorf("writing is not allowed")
			}
			return handler.AppendToFile(filePath, fileContent)
		}
	}

	return fmt.Errorf("unable to handle AppendToFile")
}

// ReadFile will read all contents of a given file
func ReadFile(filePath string) (string, error) {
	for _, handler := range Handlers {
		if handler.CanHandle(filePath) {
			if !handler.AllowRead() {
				return "", fmt.Errorf("reading is not allowed")
			}
			return handler.ReadFile(filePath)
		}
	}
	return "", fmt.Errorf("unable to handle ReadFile")
}
