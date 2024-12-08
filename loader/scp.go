package loader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/viant/afs/scp"
	"github.com/viant/afs/storage"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
	"os"
	"os/user"
	"strings"
	"time"
)

func prepareRemotePath(rpath string) (string, string, string, string, error) {
	remotePath := strings.TrimPrefix(rpath, "scp://")
	username := ""
	u, err := user.Current()
	if err == nil {
		username = u.Username
	}

	hostpath := ""
	port := "22"
	if strings.Contains(remotePath, "@") {
		parts := strings.SplitN(remotePath, "@", 2)
		username = parts[0]
		hostpath = parts[1]
	} else {
		hostpath = remotePath
	}

	if !strings.Contains(hostpath, ":") {
		return "", "", "", "", errors.New("invalid scp path, must contain <host>:<path>")
	}

	hp := strings.SplitN(hostpath, ":", 2)
	host, path := hp[0], hp[1]

	if strings.Contains(host, "*") {
		parts := strings.SplitN(host, "*", 2)
		host = parts[0]
		portPath := strings.SplitN(parts[1], ":", 2)
		port, path = portPath[0], portPath[1]
	}
	return host, path, port, username, nil
}

func getSSHConfig(username, keyPath, khPath string) (*ssh.ClientConfig, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	hkCback, err := knownhosts.New(khPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create known_hosts callback: %w", err)
	}

	return &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hkCback,
	}, nil
}

func getSCPService(host, port, key, username string) (storage.Storager, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	khPath := fmt.Sprintf("%s/.ssh/known_hosts", u.HomeDir)
	if key == "" {
		key = fmt.Sprintf("%s/.ssh/id_rsa", u.HomeDir)
	}
	conf, err := getSSHConfig(username, key, khPath)
	if err != nil {
		return nil, err
	}
	return scp.NewStorager(fmt.Sprintf("%s:%s", host, port), time.Second*10, conf)
}

// ReadFromSCP reads a file from a remote host using SCP and returns its content.
func ReadFromSCP(fileName, key string) ([]byte, string, error) {
	host, path, port, username, err := prepareRemotePath(fileName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to prepare remote path")
		return nil, "", err
	}

	service, err := getSCPService(host, port, key, username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create SCP service")
		return nil, "", err
	}

	ctx := context.Background()
	log.Debug().Msgf("Reading from path: %s", path)
	r, err := service.Open(ctx, path)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read file from %s", path)
		return nil, "", err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read file")
		return nil, "", err
	}
	return data, string(data), nil
}

// WriteToSCP writes data to a remote host using SCP.
func WriteToSCP(fileName string, data []byte, key string) error {
	host, path, port, username, err := prepareRemotePath(fileName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to prepare remote path")
		return err
	}

	service, err := getSCPService(host, port, key, username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create SCP service")
		return err
	}

	ctx := context.Background()
	d := bytes.NewReader(data)
	log.Debug().Msgf("Writing to path: %s", path)
	err = service.Upload(ctx, path, 0644, d)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to upload file to %s", path)
		return err
	}
	return nil
}
