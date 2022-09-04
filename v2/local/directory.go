package local

import "os"

func (p ProtocolHandler) ListDirectories(directoryPath string) ([]os.FileInfo, error) {
	localPath := p.ResolveFilePath(directoryPath)
	var dirs []os.FileInfo
	localDirs, err := os.ReadDir(localPath)
	if err != nil {
		return dirs, err
	}
	for _, dir := range localDirs {
		info, err := dir.Info()
		if err != nil {
			return dirs, err
		}
		dirs = append(dirs, info)
	}
	return dirs, nil
}

func (p ProtocolHandler) CreateDirectory(directoryPath string) error {
	localPath := p.ResolveFilePath(directoryPath)
	return os.MkdirAll(localPath, os.ModeDir)
}

func (p ProtocolHandler) DeleteDirectory(directoryPath string, force bool) error {
	localPath := p.ResolveFilePath(directoryPath)
	if force {
		return os.RemoveAll(localPath)
	}
	return os.Remove(localPath)
}

func (p ProtocolHandler) HasDirectory(directoryPath string) (bool, error) {
	localPath := p.ResolveFilePath(directoryPath)
	info, err := os.Stat(localPath)
	return info != nil, err
}
