package extended_fs

import (
	"fmt"

	"github.com/saitho/golang-extended-fs/core"
	"github.com/saitho/golang-extended-fs/local"
	"github.com/saitho/golang-extended-fs/pkger"
	"github.com/saitho/golang-extended-fs/sftp"
)

var Handlers = []core.HandlerFunc{
	local.ProtocolHandler{},
	sftp.ProtocolHandler{},
	pkger.ProtocolHandler{},
}

// Chown will change the owner of a file or directory
func Chown(path string, userId int, groupId int) error {
	for _, handler := range Handlers {
		if handler.CanHandle(path) {
			if !handler.AllowWrite() {
				return fmt.Errorf("writing is not allowed")
			}
			return handler.Chown(path, userId, groupId)
		}
	}
	return fmt.Errorf("unable to handle Chown")
}
