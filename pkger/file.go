package pkger

import (
	"fmt"
	"io/ioutil"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
)

func (p ProtocolHandler) ReadFile(filePath string) (string, error) {
	packagedPath := p.ResolveFilePath(filePath)
	var f pkging.File
	f, err := pkger.Open(packagedPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var sl []byte
	sl, err = ioutil.ReadAll(f)
	return string(sl), nil
}

func (p ProtocolHandler) HasFile(filePath string) (bool, error) {
	packagedPath := p.ResolveFilePath(filePath)
	var f pkging.File
	f, err := pkger.Open(packagedPath)
	if err != nil {
		return false, err
	}
	defer f.Close()
	return f != nil, err
}

func (p ProtocolHandler) WriteFile(filePath string, fileContent string) error {
	return fmt.Errorf("not supported")
}

func (p ProtocolHandler) DeleteFile(filePath string) error {
	return fmt.Errorf("not supported")
}

func (p ProtocolHandler) AppendToFile(filePath string, fileContent string) error {
	return fmt.Errorf("not supported")
}
