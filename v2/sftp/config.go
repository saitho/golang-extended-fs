package sftp

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/saitho/golang-extended-fs/v2/logger"
)

type PostDirectoryOperationHook interface {
	Execute(action DirectoryOperation, directoryName string, client *sftp.Client, data interface{}) error
}
type PostFileOperationHook interface {
	Execute(action FileOperation, fileName string, client *sftp.Client, data interface{}) error
}

type HooksStruct struct {
	PostDirectoryOperationHooks []PostDirectoryOperationHook
	PostFileOperationHooks      []PostFileOperationHook
}

type ConfigStruct struct {
	Logger logger.Wrapper
	Hooks  *HooksStruct

	SshHost                     string
	SshPort                     int
	SshUsername                 string
	SshIdentity                 string
	Signers                     []ssh.Signer
	LoadLocalSigners            bool
	AbortOnErrorsInLocalSigners bool
}

var Config = ConfigStruct{
	SshHost:                     "",
	SshPort:                     22,
	SshUsername:                 "root",
	SshIdentity:                 "",
	Signers:                     []ssh.Signer{},
	LoadLocalSigners:            true,
	AbortOnErrorsInLocalSigners: false,
}
