package core

import "os"

type TargetType int

const (
	Local TargetType = iota
	Remote
)

type HandlerFunc interface {
	GetType() TargetType
	CanHandle(filePath string) bool
	AllowRead() bool
	AllowWrite() bool
	AllowDelete() bool
	ResolveFilePath(filePath string) string

	HasDirectory(directoryPath string) (bool, error)
	ListDirectories(rootPath string) ([]os.FileInfo, error)
	CreateDirectory(directoryPath string) error
	DeleteDirectory(directoryPath string, force bool) error

	HasFile(filePath string) (bool, error)
	HasLink(filePath string) (bool, error)
	WriteFile(filePath string, fileContent string) error
	DeleteFile(filePath string) error
	AppendToFile(filePath string, fileContent string) error
	ReadFile(filePath string) (string, error)

	Chown(filePath string, userId int, groupId int) error
	Chmod(filePath string, fileMode os.FileMode) error
}
