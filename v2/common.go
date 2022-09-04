package extended_fs

import (
	"fmt"
	"os"

	"github.com/saitho/golang-extended-fs/v2/core"
	"github.com/saitho/golang-extended-fs/v2/local"
	"github.com/saitho/golang-extended-fs/v2/sftp"
)

var Handlers = []core.HandlerFunc{
	local.ProtocolHandler{},
	sftp.ProtocolHandler{},
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

// Chmod will change the permissions of a file or directory
func Chmod(path string, fileMode os.FileMode) error {
	for _, handler := range Handlers {
		if handler.CanHandle(path) {
			if !handler.AllowWrite() {
				return fmt.Errorf("writing is not allowed")
			}
			return handler.Chmod(path, fileMode)
		}
	}
	return fmt.Errorf("unable to handle Chown")
}
