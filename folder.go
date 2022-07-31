package extended_fs

import (
	"fmt"
	"os"
)

// ListFolders will list all folders within a given folder
func ListFolders(rootPath string) ([]os.FileInfo, error) {
	for _, handler := range Handlers {
		if handler.CanHandle(rootPath) {
			if !handler.AllowRead() {
				return []os.FileInfo{}, fmt.Errorf("reading is not allowed")
			}
			return handler.ListDirectories(rootPath)
		}
	}
	return []os.FileInfo{}, fmt.Errorf("unable to handle ListFolders")
}

// CreateFolder will create a new folder
func CreateFolder(folderPath string) error {
	for _, handler := range Handlers {
		if handler.CanHandle(folderPath) {
			if !handler.AllowWrite() {
				return fmt.Errorf("writing is not allowed")
			}
			return handler.CreateDirectory(folderPath)
		}
	}
	return fmt.Errorf("unable to handle CreateDirectory")
}

// DeleteFolder will delete an existing folder
func DeleteFolder(folderPath string, force bool) error {
	for _, handler := range Handlers {
		if handler.CanHandle(folderPath) {
			if !handler.AllowDelete() {
				return fmt.Errorf("deleting is not allowed")
			}
			return handler.DeleteDirectory(folderPath, force)
		}
	}
	return fmt.Errorf("unable to handle DeleteFolder")
}

// HasFolder will return true if a folder exists
func HasFolder(folderPath string) (bool, error) {
	for _, handler := range Handlers {
		if handler.CanHandle(folderPath) {
			return handler.HasDirectory(folderPath)
		}
	}
	return false, fmt.Errorf("unable to handle HasFolder")
}
