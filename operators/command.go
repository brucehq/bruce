package operators

import (
	"bruce/exe"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

type Command struct {
	Cmd        string `yaml:"cmd"`
	WorkingDir string `yaml:"dir"`
	SetEnv     string `yaml:"setEnv"`
	OnlyIf     string `yaml:"onlyIf"`
	NotIf      string `yaml:"notIf"`
	ExitIf     string `yaml:"exitIf"`
	EnvCmd     string
}

func (c *Command) Setup() {
	c.WorkingDir = RenderEnvString(c.WorkingDir)
	c.EnvCmd = RenderEnvString(c.Cmd)
	c.OnlyIf = RenderEnvString(c.OnlyIf)
	c.NotIf = RenderEnvString(c.NotIf)
	c.ExitIf = RenderEnvString(c.ExitIf)
}

// Execute runs the command.
func (c *Command) Execute() error {
	c.Setup()
	if len(c.EnvCmd) < 1 {
		return fmt.Errorf("no command to execute")
	}
	if !CanContinue(c.OnlyIf, c.NotIf, c.ExitIf, c.WorkingDir) {
		return nil
	}
	log.Info().Msgf("cmd: %s", c.EnvCmd)
	pc := exe.RunScript(c.EnvCmd, c.WorkingDir)
	if pc.Failed() {
		err := pc.GetErr()
		out := pc.Get()
		log.Error().Err(err).Msgf("failed to execute %s\n error: %s", out, err.Error())
		return err
	} else {
		log.Debug().Msgf("Output: %s", pc.Get())
		if len(c.SetEnv) > 0 {
			log.Debug().Str("cmd", c.EnvCmd).Msgf("setting env var: %s=%s", c.SetEnv, pc.Get())
			log.Error().Err(os.Setenv(c.SetEnv, pc.Get()))
		}
	}
	return nil
}
