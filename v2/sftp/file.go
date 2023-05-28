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
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "OPEN "+remotePath))
	file, err := client.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		LogError("Unable to open file at path \"" + remotePath + "\": " + err.Error())
		return err
	}
	defer file.Close()

	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "WRITE "+remotePath))
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
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "OPEN "+remotePath))
	file, err := client.OpenFile(remotePath, os.O_WRONLY|os.O_APPEND)
	if err != nil {
		LogError("Unable to open file at path \"" + remotePath + "\": " + err.Error())
		return err
	}
	defer file.Close()
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "WRITE "+remotePath))
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
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "OPEN "+remotePath))
	file, err := client.OpenFile(remotePath, os.O_RDONLY)
	if err != nil {
		LogError("Unable to open file at path \"" + remotePath + "\": " + err.Error())
		return "", err
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "READ "+remotePath))
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
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "STAT "+remotePath))
	file, err := client.Stat(remotePath)
	if err != nil && err.Error() == "file does not exist" {
		return false, nil
	}
	return file != nil, err
}

func (p ProtocolHandler) HasLink(filePath string) (bool, error) {
	remotePath := p.ResolveFilePath(filePath)
	client, err := getRemoteClient()
	if err != nil {
		return false, err
	}
	defer client.Close()
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "STAT "+remotePath))
	file, err := client.Lstat(remotePath)
	if err != nil && err.Error() == "file does not exist" {
		return false, nil
	}
	if file.Mode()&os.ModeSymlink == 0 {
		return false, fmt.Errorf("file found but it is a symlink")
	}
	return file != nil, err
}

func (p ProtocolHandler) DeleteFile(filePath string) error {
	remotePath := p.ResolveFilePath(filePath)
	client, err := getRemoteClient()
	if err != nil {
		return err
	}
	defer client.Close()
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "REMOVE "+remotePath))
	return client.Remove(remotePath)
}
