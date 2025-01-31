package operators

import (
	"bruce/loader"
	"github.com/rs/zerolog/log"
	"os"
)

type RecursiveCopy struct {
	Src           string   `yaml:"copyRecursive"`
	Dest          string   `yaml:"dest"`
	Ignores       []string `yaml:"ignoreFiles"`
	FlatCopy      bool     `yaml:"flatCopy"`
	MaxDepth      int      `yaml:"maxDepth"`
	MaxConcurrent int      `yaml:"maxConcurrent"`
	OnlyIf        string   `yaml:"onlyIf"`
	NotIf         string   `yaml:"notIf"`
	ExitIf        string   `yaml:"exitIf"`
}

func (c *RecursiveCopy) Setup() {
	c.Src = RenderEnvString(c.Src)
	c.Dest = RenderEnvString(c.Dest)
	c.OnlyIf = RenderEnvString(c.OnlyIf)
	c.NotIf = RenderEnvString(c.NotIf)
	c.ExitIf = RenderEnvString(c.ExitIf)
	// Check if parent directory exists and create it if it doesn't
	if _, err := os.Stat(c.Dest); os.IsNotExist(err) {
		err = os.MkdirAll(c.Dest, 0755)
		if err != nil {
			log.Error().Err(err).Msg("failed to create parent directory for recursive copy")
		}
	}
	if c.MaxConcurrent == 0 {
		c.MaxConcurrent = 5
	}
}

func (c *RecursiveCopy) Execute() error {
	c.Setup()
	if !CanContinue(c.OnlyIf, c.NotIf, c.ExitIf, "") {
		return nil
	}
	log.Info().Msgf("rcopy: %s => %s", c.Src, c.Dest)
	log.Info().Msgf(" (%d files at a time) with a maxDepth of: %d", c.MaxConcurrent, c.MaxDepth)

	err := loader.RecursiveCopy(c.Src, c.Dest, c.Dest, true, c.Ignores, c.FlatCopy, c.MaxDepth, c.MaxConcurrent)
	if err != nil {
		log.Error().Err(err).Msg("could not copy file")
		return err
	}
	return nil
}
