package source

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var _ Source = (*sftpSource)(nil)

type (
	sftpSource struct {
		c    *Config
		obj  io.ReadCloser
		size int64
	}
)

func openSFTPFile(c *Config) (*sftpSource, error) {
	// open connection to SFTP server
	authMethod, err := getSSHAuthMethod(c.SFTP.SSHPwd, c.SFTP.SSHKey, c.SFTP.Passphrase)
	if err != nil {
		return nil, err
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.SFTP.Host, c.SFTP.Port), &ssh.ClientConfig{
		User:            c.SFTP.SSHUser,
		Auth:            []ssh.AuthMethod{authMethod},
		Timeout:         time.Second * 5,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint: gosec
	})
	if err != nil {
		return nil, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	// open the file
	obj, err := client.Open(c.SFTP.Path)
	if err != nil {
		return nil, err
	}

	// get the file size
	stat, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	return &sftpSource{
		c:    c,
		obj:  obj,
		size: stat.Size(),
	}, nil
}

func (s *sftpSource) Config() *Config {
	return s.c
}

func (s *sftpSource) Size() (int64, error) {
	return s.size, nil
}

func (s *sftpSource) Read(p []byte) (int, error) {
	return s.obj.Read(p)
}

func (s *sftpSource) Close() error {
	return s.obj.Close()
}

func getSSHAuthMethod(sshPwd, sshKey, passphrase string) (ssh.AuthMethod, error) {
	if sshKey != "" {
		key, err := getSSHSigner(sshKey, passphrase)
		if err != nil {
			return nil, err
		}
		return ssh.PublicKeys(key), nil
	}
	return ssh.Password(sshPwd), nil
}

func getSSHSigner(keyData, passphrase string) (ssh.Signer, error) {
	if passphrase != "" {
		return ssh.ParsePrivateKeyWithPassphrase([]byte(keyData), []byte(passphrase))
	}
	return ssh.ParsePrivateKey([]byte(keyData))
}
