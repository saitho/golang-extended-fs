package sftp

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
	return core.Remote
}
func (p ProtocolHandler) CanHandle(filePath string) bool {
	return strings.HasPrefix(filePath, "sftp://") || strings.HasPrefix(filePath, "ssh://")
}

func (p ProtocolHandler) ResolveFilePath(filePath string) string {
	if strings.HasPrefix(filePath, "sftp://") {
		return filePath[7:] // sftp://
	}
	return filePath[6:] // ssh://
}

func (p ProtocolHandler) Chown(directoryPath string, userId int, groupId int) error {
	remotePath := p.ResolveFilePath(directoryPath)
	client, err := getRemoteClient()
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Chown(remotePath, userId, groupId)
}
