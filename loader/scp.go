package loader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/viant/afs/scp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
	"os"
	"os/user"
	"strings"
	"time"
)

func prepareRemotePath(rpath string) (string, string, string, string, error) {
	// first we strip off the first 6 chars which correspond to scp://
	remotePath := rpath[6:]
	u, err := user.Current()
	username := ""
	if err == nil {
		username = u.Username
	}

	hostpath := ""
	port := "22"
	// now we look to see if it's using a user name before the @ and assign the rest to "hostpath" variable
	if strings.Contains(remotePath, "@") {
		// split the remote path into two parts, the user and the hostpath
		username = strings.Split(remotePath, "@")[0]
		hostpath = strings.Split(remotePath, "@")[1]
	} else {
		hostpath = remotePath
	}
	if !strings.Contains(remotePath, ":") {
		log.Debug().Msgf("invalid scp path, must contain remote path separated by : %s", remotePath)
		return "", "", "", "", errors.New("invalid scp path, must contain <host>:<path>")
	}
	log.Debug().Msgf("host path: %s", hostpath)
	hp := strings.Split(hostpath, ":")
	if len(hp) < 2 {
		log.Debug().Msgf("Current path pieces: %#v", hp)
		return "", "", "", "", fmt.Errorf("invalid scp path: %s", hostpath)
	}
	host := hp[0]
	path := hp[1]
	if strings.Contains(host, "*") {
		pSplit := strings.Split(host, "*")
		host = pSplit[0]
		portPath := strings.Split(pSplit[1], ":")
		port = portPath[0]
		path = portPath[1]
	}
	return host, path, port, username, nil
}

// ReadFromSCP will use RSSH to create a connection then read the file from the remote host and save it to a local file with io.Copy
func ReadFromSCP(fileName, key string) ([]byte, string, error) {
	host, path, port, username, err := prepareRemotePath(fileName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to prepare remote path")
		return nil, "", err
	}
	u, kerr := user.Current()
	if kerr != nil {
		log.Error().Err(err).Msg("Failed to get current user")
		return nil, "", err
	}
	khPath := fmt.Sprintf("%s/.ssh/known_hosts", u.HomeDir)
	if key == "" {
		key = fmt.Sprintf("%s/.ssh/id_rsa", u.HomeDir)
	}

	log.Debug().Msgf("reading key file: %s", key)
	keyData, err := os.ReadFile(key)
	if err != nil {
		log.Error().Err(err).Msg("failed to read private key")
		return nil, "", fmt.Errorf("failed to read private key: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse private key: %w", err)
	}
	hkCback, err := knownhosts.New(khPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create known_hosts callback: %w", err)
	}

	conf := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hkCback,
	}
	service, err := scp.NewStorager(fmt.Sprintf("%s:%s", host, port), time.Duration(time.Second*10), conf)
	if err != nil {
		log.Error().Err(err).Msg("failed to create scp service")
		return nil, "", err
	}
	ctx := context.Background()
	log.Debug().Msgf("path to read: %s", path)
	r, err := service.Open(ctx, path)
	if err != nil {
		log.Error().Err(err).Msgf("failed to upload file to %s", path)
		return nil, "", err
	}
	d, err := io.ReadAll(r)
	if err != nil {
		log.Error().Err(err).Msg("failed to read file")
		return nil, "", err
	}
	return d, string(d), nil
}

// WriteToSCP will use RSSH to create a connection then write the file to the remote host with io.Copy
func WriteToSCP(fileName string, data []byte, key string) error {
	host, path, port, username, err := prepareRemotePath(fileName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to prepare remote path")
		return err
	}
	u, kerr := user.Current()
	if kerr != nil {
		log.Error().Err(err).Msg("Failed to get current user")
		return err
	}
	khPath := fmt.Sprintf("%s/.ssh/known_hosts", u.HomeDir)
	if key == "" {
		key = fmt.Sprintf("%s/.ssh/id_rsa", u.HomeDir)
	}

	log.Debug().Msgf("reading key file: %s", key)
	keyData, err := os.ReadFile(key)
	if err != nil {
		log.Error().Err(err).Msg("failed to read private key")
		return fmt.Errorf("failed to read private key: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	hkCback, err := knownhosts.New(khPath)
	if err != nil {
		return fmt.Errorf("failed to create known_hosts callback: %w", err)
	}

	conf := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hkCback,
	}
	service, err := scp.NewStorager(fmt.Sprintf("%s:%s", host, port), time.Duration(time.Second*10), conf)
	if err != nil {
		log.Error().Err(err).Msg("failed to create scp service")
		return err
	}
	ctx := context.Background()
	d := bytes.NewReader(data)
	log.Debug().Msgf("path to write: %s", path)
	err = service.Upload(ctx, path, 0644, d)
	if err != nil {
		log.Error().Err(err).Msgf("failed to upload file to %s", path)
		return err
	}
	return nil
}
