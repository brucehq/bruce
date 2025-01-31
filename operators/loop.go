package operators

import (
	"bruce/exe"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

type Loop struct {
	LoopScript string `yaml:"loopScript"`
	Count      int    `yaml:"count"`
	Variable   string `yaml:"var"`
	OnlyIf     string `yaml:"onlyIf"`
	NotIf      string `yaml:"notIf"`
	ExitIf     string `yaml:"exitIf"`
}

func (lp *Loop) Setup() {
	lp.OnlyIf = RenderEnvString(lp.OnlyIf)
	lp.NotIf = RenderEnvString(lp.NotIf)
	lp.ExitIf = RenderEnvString(lp.ExitIf)
	lp.LoopScript = RenderEnvString(lp.LoopScript)
}

// Execute runs the command.
func (lp *Loop) Execute() error {
	lp.Setup()
	/* We do not replace command envars like the other functions, this is intended to be a raw command */
	// if onlyIf is set, check if it's return value is not empty / true
	if !CanContinue(lp.OnlyIf, lp.NotIf, lp.ExitIf, "") {
		return nil
	}
	for i := 0; i < lp.Count; i++ {
		log.Info().Str("loop", lp.LoopScript).Msgf("executing: %s with variable: %s and value: %d", lp.LoopScript, lp.Variable, i)
		log.Error().Err(os.Setenv(lp.Variable, fmt.Sprintf("%d", i)))
		// get current running file and append the loop script as the first argument
		execCmd := fmt.Sprintf("%s %s", os.Args[0], lp.LoopScript)
		pc := exe.Run(execCmd, "")
		if pc.Failed() {
			log.Error().Err(pc.GetErr()).Msg(pc.Get())
			return pc.GetErr()
		}
	}

	return nil
}
