package operators

import (
	"bruce/mutation"
	"fmt"
	"github.com/rs/zerolog/log"
)

type Tarball struct {
	Name   string `yaml:"name"`
	Src    string `yaml:"tarball"`
	Dest   string `yaml:"dest"`
	Force  bool   `yaml:"force"`
	Strip  bool   `yaml:"stripRoot"`
	OnlyIf string `yaml:"onlyIf"`
	NotIf  string `yaml:"notIf"`
	ExitIf string `yaml:"exitIf"`
}

func (t *Tarball) Setup() {
	t.Src = RenderEnvString(t.Src)
	t.Dest = RenderEnvString(t.Dest)
	t.OnlyIf = RenderEnvString(t.OnlyIf)
	t.NotIf = RenderEnvString(t.NotIf)
	t.ExitIf = RenderEnvString(t.ExitIf)
}

func (t *Tarball) Execute() error {
	t.Setup()
	if !CanContinue(t.OnlyIf, t.NotIf, t.ExitIf, "") {
		return nil
	}
	log.Info().Msgf("tarball: %s => %s", t.Src, t.Dest)
	if len(t.Src) < 1 {
		return fmt.Errorf("source is too short")
	}
	log.Info().Msgf("tarball: %s => %s", t.Src, t.Dest)
	return mutation.ExtractTarball(t.Src, t.Dest, t.Force, t.Strip)
}
