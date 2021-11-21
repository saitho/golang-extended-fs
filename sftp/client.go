package sftp

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type ClientConfigHook interface {
	CreateClientConfig(*ssh.ClientConfig) *ssh.ClientConfig
}

var clientConfigHooks []ClientConfigHook
func RegisterClientConfigHook(hookFunc ClientConfigHook) {
	clientConfigHooks = append(clientConfigHooks, hookFunc)
}

var SIZE = 1 << 15

func getRemoteSshAuths() ([]ssh.AuthMethod, error) {
	var auths []ssh.AuthMethod
	var signers = []ssh.Signer{}
	signers = append(signers, Config.Signers...)

	if Config.LoadLocalSigners {
		if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
			localSigners, err := agent.NewClient(aconn).Signers()
			if err != nil {
				if Config.AbortOnErrorsInLocalSigners {
					return auths, err
				}
			} else {
				signers = append(signers, localSigners...)
			}
		}
	}

	auths = append(auths, ssh.PublicKeys(signers...))
	return auths, nil
}

func getRemoteClient() (*sftp.Client, error) {
	auths, err := getRemoteSshAuths()
	if err != nil {
		return nil, err
	}

	clientConfig := &ssh.ClientConfig{
		User:            Config.SshUsername,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	for _, hook := range clientConfigHooks {
		clientConfig = hook.CreateClientConfig(clientConfig)
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", Config.SshHost, Config.SshPort), clientConfig)
	if err != nil {
		return nil, err
	}

	return sftp.NewClient(conn, sftp.MaxPacket(SIZE))
}

func RemoteRun(cmd string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	remoteCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	if Config.Logger != nil {
		(*Config.Logger).Debug(fmt.Sprintf("SSH [%s@%s]: %s", Config.SshUsername, Config.SshHost, remoteCmd))
	}

	var cmdArgs []string
	if Config.SshIdentity != "" {
		cmdArgs = []string{"-i", Config.SshIdentity}
	}

	cmdArgs = append(cmdArgs, fmt.Sprintf("%s@%s", Config.SshUsername, Config.SshHost))
	cmdArgs = append(cmdArgs, remoteCmd)
	command := exec.Command("ssh", cmdArgs...)
	var out, outErr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &outErr
	err := command.Run()
	return out, outErr, err
}
