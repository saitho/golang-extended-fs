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

// DeleteFile will delete the given file
func DeleteFile(filePath string) error {
	for _, handler := range Handlers {
		if handler.CanHandle(filePath) {
			if !handler.AllowDelete() {
				return fmt.Errorf("deleting is not allowed")
			}
			return handler.DeleteFile(filePath)
		}
	}
	return fmt.Errorf("unable to handle DeleteFile")
}

// CopyFile will read all contents of a given file and write it to another location
func CopyFile(srcPath string, destPath string) error {
	content, err := ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("unable to read source file")
	}
	if err := WriteFile(destPath, content); err != nil {
		return fmt.Errorf("unable to write to source file")
	}
	return nil
}
