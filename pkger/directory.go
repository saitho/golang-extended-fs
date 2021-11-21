package pkger

import (
	"fmt"
	"os"
)

func (p ProtocolHandler) ListDirectories(directoryPath string) ([]os.FileInfo, error) {
	return []os.FileInfo{}, fmt.Errorf("not supported")
}

func (p ProtocolHandler) CreateDirectory(directoryPath string) error {
	return fmt.Errorf("not supported")
}
