package ssh_server

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type SshClient struct {
	IP         string      //IP地址
	Port       int         //端口号
	client     *ssh.Client //ssh客户端
	config     *ssh.ClientConfig
	LastResult string //最近一次Run的结果
}

func Init(ip string, username string, password string) *SshClient {
	cli := new(SshClient)
	cli.IP = ip
	cli.config = &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
		Timeout:         30 * time.Second,
	}
	return cli
}

func (c *SshClient) Run(shell string) (string, error) {
	if c.client == nil {
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", c.IP), c.config)
		if err != nil {
			return "", err
		}
		c.client = client
	}
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	buf, err := session.CombinedOutput(shell)
	c.LastResult = string(buf)
	return c.LastResult, err
}
