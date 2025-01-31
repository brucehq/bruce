package operators

import (
	"bruce/mutation"
	"bruce/system"
	"fmt"
	"github.com/rs/zerolog/log"
	"runtime"
)

// Cron provides a means to set the ownership of files or directories as needed.
type Cron struct {
	Name     string `yaml:"cron"`
	Schedule string `yaml:"schedule"`
	User     string `yaml:"username"`
	Exec     string `yaml:"cmd"`
	OnlyIf   string `yaml:"onlyIf"`
	NotIf    string `yaml:"notIf"`
	ExitIf   string `yaml:"exitIf"`
}

func (c *Cron) Setup() {
	c.Exec = RenderEnvString(c.Exec)
	c.User = RenderEnvString(c.User)
	c.OnlyIf = RenderEnvString(c.OnlyIf)
	c.NotIf = RenderEnvString(c.NotIf)
	c.ExitIf = RenderEnvString(c.ExitIf)
}

func (c *Cron) Execute() error {
	c.Setup()
	if runtime.GOOS == "linux" {
		if !CanContinue(c.OnlyIf, c.NotIf, c.ExitIf, "") {
			return nil
		}
		jobName := mutation.StripNonAlnum(c.Name)
		log.Info().Msgf("cron: /etc/cron.d/%s", jobName)
		c.Schedule = mutation.StripExtraWhitespaceFB(c.Schedule)
		c.User = mutation.StripNonAlnum(c.User)
		log.Debug().Msgf("starting cronjob: %s", jobName)
		if c.User == "" {
			c.User = system.Get().CurrentUser.Username
		}
		return mutation.WriteInlineTemplate(fmt.Sprintf("/etc/cron.d/%s", jobName), "{{.Schedule}} {{.User}} {{.Exec}}", c)
	}
	return fmt.Errorf("not supported")
}
