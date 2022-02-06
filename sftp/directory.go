package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"os"
)

type DirectoryOperation int

const (
	DirList   DirectoryOperation = iota
	DirCreate DirectoryOperation = 1
)

func triggerDirectoryHooks(operation DirectoryOperation, remotePath string, client *sftp.Client, data interface{}) error {
	if Config.Hooks != nil && Config.Hooks.PostDirectoryOperationHooks != nil {
		for _, hook := range Config.Hooks.PostDirectoryOperationHooks {
			if err := hook.Execute(operation, remotePath, client, data); err != nil {
				LogError(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "PostDirectoryOperationHook (operation: "+string(rune(operation))+" failed on  "+remotePath+": "+err.Error()))
				return err
			}
		}
	}
	return nil
}

func (p ProtocolHandler) ListDirectories(rootPath string) ([]os.FileInfo, error) {
	remotePath := p.ResolveFilePath(rootPath)
	var dirs []os.FileInfo
	client, err := getRemoteClient()
	if err != nil {
		return dirs, err
	}
	defer client.Close()
	files, err := client.ReadDir(remotePath)
	if err != nil {
		LogError("Unable to open directory at path \"" + remotePath + "\": " + err.Error())
		return dirs, err
	}
	if err := triggerDirectoryHooks(DirList, remotePath, client, files); err != nil {
		return files, err
	}
	return files, nil
}

func (p ProtocolHandler) CreateDirectory(directoryPath string) error {
	remotePath := p.ResolveFilePath(directoryPath)
	client, err := getRemoteClient()
	if err != nil {
		return err
	}
	defer client.Close()
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "MKDIRALL "+remotePath))
	if err := client.MkdirAll(remotePath); err != nil {
		LogError(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "MKDIRALL "+remotePath+" failed: "+err.Error()))
		return err
	}

	return triggerDirectoryHooks(DirCreate, remotePath, client, nil)
}
