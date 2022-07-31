package local

import (
	"os"
	"strings"

	"github.com/saitho/golang-extended-fs/core"
)

type ProtocolHandler struct {
}

func (p ProtocolHandler) AllowRead() bool {
	return true
}
func (p ProtocolHandler) AllowWrite() bool {
	return true
}
func (p ProtocolHandler) AllowDelete() bool {
	return true
}

func (p ProtocolHandler) GetType() core.TargetType {
	return core.Local
}

func (p ProtocolHandler) CanHandle(filePath string) bool {
	return !strings.Contains(filePath, "://")
}

func (p ProtocolHandler) ResolveFilePath(filePath string) string {
	return filePath
}

func (p ProtocolHandler) Chown(directoryPath string, userId int, groupId int) error {
	remotePath := p.ResolveFilePath(directoryPath)
	return os.Chown(remotePath, userId, groupId)
}
