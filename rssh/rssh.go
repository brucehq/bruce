package rssh

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/knownhosts"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type RSSH struct {
	Host   string
	User   string
	Key    []byte
	Port   string
	client *ssh.Client
}

func NewRSSH(host, user, privkey string, allowInsecure bool) (*RSSH, error) {
	if privkey == "" {
		privkey = os.ExpandEnv("$HOME/.ssh/id_rsa")
	}
	keyBytes, err := os.ReadFile(privkey)
	if err != nil {
		log.Error().Err(err).Msg("failed to read private key")
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	h, p := formatHostAndPort(host)
	rsshc := &RSSH{
		Host: h, // Ensure IPv6 addresses are correctly formatted
		User: user,
		Key:  keyBytes,
		Port: p,
	}
	err = rsshc.setup(allowInsecure)
	if err != nil {
		return nil, err
	}
	return rsshc, nil
}

func formatHostAndPort(hostport string) (string, string) {
	// parse, detect IPv4 vs IPv6, return correct host + port
	// pseudocode:
	if strings.Contains(hostport, ":") {
		splits := strings.Split(hostport, ":")
		host := splits[0]
		port := splits[1]
		// naive check for IPv6
		if strings.Contains(host, ":") {
			return "[" + host + "]", ":" + port
		}
		return host, port
	}
	// default to port 22
	return hostport, "22"
}

func (r *RSSH) setup(doInsecure bool) error {

	signer, err := ssh.ParsePrivateKey(r.Key)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	khPath := os.ExpandEnv("$HOME/.ssh/known_hosts")
	hkCback, err := knownhosts.New(khPath)
	if err != nil {
		return fmt.Errorf("failed to create known_hosts callback: %w", err)
	}
	iCBack := ssh.InsecureIgnoreHostKey()

	config := &ssh.ClientConfig{
		User: r.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hkCback,
	}
	if doInsecure {
		config.HostKeyCallback = iCBack
	}
	addr, port := formatHostAndPort(r.Host)
	client, err := ssh.Dial("tcp", net.JoinHostPort(addr, port), config)
	if err != nil {
		return err
	}
	r.client = client
	return nil
}

func (r *RSSH) ExecCommand(cmd string) (string, error) {
	session, err := r.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &b

	err = session.Run(cmd)
	if err != nil {
		log.Debug().Err(err).Msg("remote command execution failed")
		return "", err
	}
	return b.String(), nil
}

func (r *RSSH) Close() {
	if r.client != nil {
		log.Error().Err(r.client.Close())
		r.client = nil
	}
}
