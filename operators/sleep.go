package operators

import (
	"github.com/rs/zerolog/log"
	"time"
)

type Sleep struct {
	Time   int    `yaml:"sleep"`
	OnlyIf string `yaml:"onlyIf"`
	NotIf  string `yaml:"notIf"`
	ExitIf string `yaml:"exitIf"`
}

func (sl *Sleep) Setup() {
	sl.OnlyIf = RenderEnvString(sl.OnlyIf)
	sl.NotIf = RenderEnvString(sl.NotIf)
	sl.ExitIf = RenderEnvString(sl.ExitIf)
}

// Execute runs the command.
func (sl *Sleep) Execute() error {
	sl.Setup()
	if !CanContinue(sl.OnlyIf, sl.NotIf, sl.ExitIf, "") {
		return nil
	}
	log.Info().Msgf("sleep: %d seconds", sl.Time)
	time.Sleep(time.Duration(sl.Time) * time.Second)

	return nil
}
