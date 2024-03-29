package sftp

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pkg/sftp"
)

type FileOperation int

const (
	FileCreate FileOperation = iota
)

func triggerFileHooks(operation FileOperation, remotePath string, client *sftp.Client, data interface{}) error {
	if Config.Hooks != nil && Config.Hooks.PostFileOperationHooks != nil {
		for _, hook := range Config.Hooks.PostFileOperationHooks {
			if err := hook.Execute(operation, remotePath, client, data); err != nil {
				LogError(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "PostFileOperationHooks (operation: "+string(rune(operation))+" failed on  "+remotePath+": "+err.Error()))
				return err
			}
		}
	}
	return nil
}

func (p ProtocolHandler) WriteFile(filePath string, fileContent string) error {
	remotePath := p.ResolveFilePath(filePath)
	client, err := getRemoteClient()
	if err != nil {
		return err
	}
	defer client.Close()
	file, err := client.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		LogError("Unable to open file at path \"" + remotePath + "\": " + err.Error())
		return err
	}
	defer file.Close()

	if _, err := file.Write([]byte(fileContent)); err != nil {
		LogError("Unable to write file at path \"" + remotePath + "\": " + err.Error())
		return err
	}
	return triggerFileHooks(FileCreate, remotePath, client, file)
}

func (p ProtocolHandler) AppendToFile(filePath string, fileContent string) error {
	remotePath := p.ResolveFilePath(filePath)
	client, err := getRemoteClient()
	if err != nil {
		return err
	}
	defer client.Close()
	file, err := client.OpenFile(remotePath, os.O_WRONLY|os.O_APPEND)
	if err != nil {
		LogError("Unable to open file at path \"" + remotePath + "\": " + err.Error())
		return err
	}
	defer file.Close()
	if _, err := file.Write([]byte(fileContent)); err != nil {
		LogError("Unable to write file at path \"" + remotePath + "\": " + err.Error())
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
		LogError("Unable to open file at path \"" + remotePath + "\": " + err.Error())
		return "", err
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		LogError("Unable to read file at path \"" + remotePath + "\": " + err.Error())
		return "", err
	}
	return buf.String(), nil
}

func (p ProtocolHandler) HasFile(filePath string) (bool, error) {
	remotePath := p.ResolveFilePath(filePath)
	client, err := getRemoteClient()
	if err != nil {
		return false, err
	}
	defer client.Close()
	file, err := client.Stat(remotePath)
	return file != nil, err
}

func (p ProtocolHandler) DeleteFile(filePath string) error {
	remotePath := p.ResolveFilePath(filePath)
	client, err := getRemoteClient()
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Remove(remotePath)
}
