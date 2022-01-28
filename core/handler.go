package core

import "os"

type TargetType int

const (
	Pkging TargetType = iota
	Local
	Remote
)

type HandlerFunc interface {
	GetType() TargetType
	CanHandle(filePath string) bool
	AllowRead() bool
	AllowWrite() bool
	ResolveFilePath(filePath string) string

	ListDirectories(rootPath string) ([]os.FileInfo, error)
	CreateDirectory(directoryPath string) error
	WriteFile(filePath string, fileContent string) error
	AppendToFile(filePath string, fileContent string) error
	ReadFile(filePath string) (string, error)
	Chown(filePath string, userId int, groupId int) error
}
