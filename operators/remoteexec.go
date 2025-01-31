package operators

import (
	"bruce/rssh"
	"github.com/rs/zerolog/log"
	"os"
	"os/user"
	"strings"
)

type RemoteExec struct {
	ExecCmd       string `yaml:"remoteCmd"`
	RemHost       string `yaml:"host"`
	SetEnv        string `yaml:"setEnv"`
	PrivKey       string `yaml:"key"`
	AllowInsecure bool   `yaml:"allowInsecure"`
	OnlyIf        string `yaml:"onlyIf"`
	NotIf         string `yaml:"notIf"`
	ExitIf        string `yaml:"exitIf"`
}

func (re *RemoteExec) Setup() {
	re.ExecCmd = RenderEnvString(re.ExecCmd)
	re.RemHost = RenderEnvString(re.RemHost)
	re.OnlyIf = RenderEnvString(re.OnlyIf)
	re.NotIf = RenderEnvString(re.NotIf)
	re.ExitIf = RenderEnvString(re.ExitIf)
}

func (re *RemoteExec) Execute() error {
	re.Setup()
	if !CanContinue(re.OnlyIf, re.NotIf, re.ExitIf, "") {
		return nil
	}
	log.Info().Msgf("remoteCmd: (%s): %s", re.RemHost, re.ExecCmd)
	usr, err := user.Current()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current user")
		return err
	}
	uname := usr.Username
	hostname := re.RemHost
	if strings.Contains(re.RemHost, "@") {
		uname = strings.Split(re.RemHost, "@")[0]
		hostname = strings.Split(re.RemHost, "@")[1]
	}
	rs, err := rssh.NewRSSH(hostname, uname, re.PrivKey, re.AllowInsecure)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create RSSH")
		return err
	}
	defer rs.Close()
	output, err := rs.ExecCommand(re.ExecCmd)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to execute %s", re.ExecCmd)
		return err
	}
	log.Debug().Str("cmd", re.ExecCmd).Msgf("completed executing on [%s]", re.RemHost)
	log.Debug().Msgf("Output: %s", output)
	if len(re.SetEnv) > 0 {
		log.Debug().Str("remoteCmd", re.SetEnv).Msgf("setting env var: %s=%s", re.SetEnv, output)
		log.Error().Err(os.Setenv(re.SetEnv, output))
	}

	return nil
}
