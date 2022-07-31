package pkger

import (
	"fmt"
	"strings"

	"github.com/saitho/golang-extended-fs/core"
)

type ProtocolHandler struct {
}

func (p ProtocolHandler) AllowRead() bool {
	return true
}
func (p ProtocolHandler) AllowWrite() bool {
	return false
}
func (p ProtocolHandler) AllowDelete() bool {
	return false
}

func (p ProtocolHandler) GetType() core.TargetType {
	return core.Pkging
}
func (p ProtocolHandler) CanHandle(filePath string) bool {
	return strings.HasPrefix(filePath, "pkging://")
}

func (p ProtocolHandler) ResolveFilePath(filePath string) string {
	return filePath[9:]
}

func (p ProtocolHandler) Chown(directoryPath string, userId int, groupId int) error {
	return fmt.Errorf("not supported")
}
