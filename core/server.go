package core

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

type Server struct {
	Name     string `json:"name"`
	Method   string `json:"method"`
	User     string `json:"user"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

func (s *Server) Connect() {
	auth, err := s.parseAuthMethod()
	if err != nil {
		panic(err)
		return
	}

	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{auth},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := s.IP + ":" + strconv.Itoa(s.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		panic(err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic(err)
		return
	}

	defer session.Close()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
		return
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
		return
	}

	defer terminal.Restore(fd, state)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return
	}

	err = session.Shell()
	if err != nil {
		panic(err)
		return
	}

	err = session.Wait()
	if err != nil {
		panic(err)
		return
	}
}

func (s *Server) publicKeyFile() (ssh.AuthMethod, error) {
	pemBytes, err := ioutil.ReadFile(s.Key)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if s.Password == "" {
		signer, err = ssh.ParsePrivateKey(pemBytes)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(s.Password))
	}

	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}

func (s *Server) parseAuthMethod() (ssh.AuthMethod, error) {
	var m ssh.AuthMethod

	switch s.Method {
	case "password":
		m = ssh.Password(s.Password)
		break

	case "pem":
		mtd, err := s.publicKeyFile()
		if err != nil {
			return nil, err
		}
		m = mtd
		break

	default:
		return nil, errors.New("密码类型错误: " + s.Method)
	}

	return m, nil
}