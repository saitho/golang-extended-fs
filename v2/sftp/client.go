package sftp

import (
	"fmt"
	"net"
	"os"
	"strconv"

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

	if Config.SshIdentity != "" {
		dat, err := os.ReadFile(Config.SshIdentity)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(dat)
			if err == nil {
				signers = append(signers, signer)
			}
		}
	}

	auths = append(auths, ssh.PublicKeys(signers...))
	return auths, nil
}

func getRemoteClient() (*sftp.Client, error) {
	auths, err := getRemoteSshAuths()
	if err != nil {
		LogError("Unable to get remote ssh auths: " + err.Error())
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

	addr := net.JoinHostPort(Config.SshHost, strconv.Itoa(Config.SshPort))
	LogDebug("Dialing " + addr + " with user " + Config.SshUsername)
	conn, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		LogError("Unable to connect to host " + fmt.Sprintf("%s:%d", Config.SshHost, Config.SshPort) + ": " + err.Error())
		return nil, err
	}

	return sftp.NewClient(conn, sftp.MaxPacket(SIZE))
}
