package sftp

import (
	"fmt"
	"os"
)

type DirectoryOperation int

const (
	DirList   DirectoryOperation = iota
	DirCreate DirectoryOperation = 1
)

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
		return dirs, err
	}
	if Config.Hooks != nil && Config.Hooks.PostDirectoryOperationHooks != nil {
		for _, hook := range Config.Hooks.PostDirectoryOperationHooks {
			if err := hook.Execute(DirList, remotePath, client, files); err != nil {
				return files, err
			}
		}
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
	if Config.Logger != nil {
		(*Config.Logger).Debug(fmt.Sprintf("SFTP [%s@%s]: %s", Config.SshUsername, Config.SshHost, "MKDIRALL "+remotePath))
	}
	if err := client.MkdirAll(remotePath); err != nil {
		return err
	}

	if Config.Hooks != nil && Config.Hooks.PostDirectoryOperationHooks != nil {
		for _, hook := range Config.Hooks.PostDirectoryOperationHooks {
			if err := hook.Execute(DirCreate, remotePath, client, nil); err != nil {
				return err
			}
		}
	}

	return nil
}
