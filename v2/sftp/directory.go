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
	DirRemove DirectoryOperation = 2
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
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "READ "+remotePath))
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

func deleteRecursively(client *sftp.Client, path string) error {
	fileInfo, err := client.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range fileInfo {
		fullPath := path + "/" + file.Name()
		if file.IsDir() {
			if err := deleteRecursively(client, fullPath); err != nil {
				return err
			}
		}
		if err := client.Remove(fullPath); err != nil {
			return err
		}
	}
	return nil
}

func (p ProtocolHandler) DeleteDirectory(directoryPath string, force bool) error {
	remotePath := p.ResolveFilePath(directoryPath)
	client, err := getRemoteClient()
	if err != nil {
		return err
	}
	defer client.Close()
	if force {
		// delete folder contents
		if err := deleteRecursively(client, remotePath); err != nil {
			return err
		}
	}
	LogDebug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "REMOVE "+remotePath))
	if err := client.Remove(remotePath); err != nil {
		LogError(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "REMOVE "+remotePath+" failed: "+err.Error()))
		return err
	}
	return triggerDirectoryHooks(DirRemove, remotePath, client, nil)
}

func (p ProtocolHandler) HasDirectory(directoryPath string) (bool, error) {
	remotePath := p.ResolveFilePath(directoryPath)
	client, err := getRemoteClient()
	if err != nil {
		return false, err
	}
	defer client.Close()
	info, err := client.Stat(remotePath)
	if err != nil && err.Error() == "file does not exist" {
		return false, nil
	}
	return info != nil, err
}
