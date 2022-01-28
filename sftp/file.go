package sftp

import (
	"bytes"
	"os"
)

type FileOperation int

const (
	FileCreate FileOperation = iota
)

func (p ProtocolHandler) WriteFile(filePath string, fileContent string) error {
	remotePath := p.ResolveFilePath(filePath)
	client, err := getRemoteClient()
	if err != nil {
		return err
	}
	defer client.Close()
	file, err := client.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write([]byte(fileContent)); err != nil {
		return err
	}
	for _, hook := range Config.Hooks.PostFileOperationHooks {
		if err := hook.Execute(FileCreate, remotePath, client, file); err != nil {
			return err
		}
	}
	return nil
}

func (p ProtocolHandler) AppendToFile(remotePath string, fileContent string) error {
	client, err := getRemoteClient()
	if err != nil {
		return err
	}
	defer client.Close()
	file, err := client.OpenFile(remotePath, os.O_WRONLY|os.O_APPEND)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write([]byte(fileContent)); err != nil {
		return err
	}
	return nil
}

func (p ProtocolHandler) ReadFile(filePath string) (string, error) {
	remotePath := p.ResolveFilePath(filePath)
	client, err := getRemoteClient()
	if err != nil {
		return "", err
	}
	defer client.Close()
	file, err := client.OpenFile(remotePath, os.O_RDONLY)
	if err != nil {
		return "", err
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", err
	}
	return buf.String(), nil
}
