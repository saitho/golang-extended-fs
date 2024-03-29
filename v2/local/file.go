package local

import (
	"fmt"
	"os"
)

func (p ProtocolHandler) WriteFile(filePath string, fileContent string) error {
	localPath := p.ResolveFilePath(filePath)
	f, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer f.Close()
	if err != nil {
		return err
	}
	if _, err := f.WriteString(fileContent); err != nil {
		return err
	}
	return nil
}

func (p ProtocolHandler) AppendToFile(filePath string, fileContent string) error {
	localPath := p.ResolveFilePath(filePath)
	f, err := os.OpenFile(localPath, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(fileContent); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func (p ProtocolHandler) ReadFile(filePath string) (string, error) {
	localPath := p.ResolveFilePath(filePath)
	dat, err := os.ReadFile(localPath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func (p ProtocolHandler) DeleteFile(filePath string) error {
	localPath := p.ResolveFilePath(filePath)
	return os.Remove(localPath)
}

func (p ProtocolHandler) HasFile(filePath string) (bool, error) {
	localPath := p.ResolveFilePath(filePath)
	info, err := os.Stat(localPath)
	if err != nil && err.Error() == "file does not exist" {
		return false, nil
	}
	return info != nil, err
}

func (p ProtocolHandler) HasLink(filePath string) (bool, error) {
	localPath := p.ResolveFilePath(filePath)
	info, err := os.Lstat(localPath)
	if err != nil && err.Error() == "link does not exist" {
		return false, nil
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return false, fmt.Errorf("file found but it is a symlink")
	}
	return info != nil, err
}
