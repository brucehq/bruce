package operators

import (
	"bruce/loader"
	"github.com/rs/zerolog/log"
	"io/fs"
)

type Copy struct {
	Src    string      `yaml:"copy"`
	Dest   string      `yaml:"dest"`
	Key    string      `yaml:"key"`
	Perm   fs.FileMode `yaml:"perm"`
	OnlyIf string      `yaml:"onlyIf"`
	NotIf  string      `yaml:"notIf"`
	ExitIf string      `yaml:"exitIf"`
}

func (c *Copy) Setup() {
	c.Src = RenderEnvString(c.Src)
	c.Dest = RenderEnvString(c.Dest)
	c.OnlyIf = RenderEnvString(c.OnlyIf)
	c.NotIf = RenderEnvString(c.NotIf)
	c.ExitIf = RenderEnvString(c.ExitIf)
}

func (c *Copy) Execute() error {
	c.Setup()
	if !CanContinue(c.OnlyIf, c.NotIf, c.ExitIf, "") {
		return nil
	}
	log.Info().Msgf("copy: %s => %s", c.Src, c.Dest)
	err := loader.CopyFile(c.Src, c.Dest, c.Key, c.Perm, true)
	log.Info().Msgf("copy: %s => %s", c.Src, c.Dest)
	if err != nil {
		log.Error().Err(err).Msg("could not copy file")
		return err
	}
	return nil
}
