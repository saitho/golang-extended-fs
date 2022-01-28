package local

import (
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

func (p ProtocolHandler) GetType() core.TargetType {
	return core.Local
}

func (p ProtocolHandler) CanHandle(filePath string) bool {
	return !strings.Contains(filePath, "://")
}

func (p ProtocolHandler) ResolveFilePath(filePath string) string {
	return filePath
}
