package operators

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"os"
	"path"
)

type Git struct {
	Repo     string `yaml:"gitRepo"`
	Location string `yaml:"dest"`
	Mode     string `yaml:"mode"` // cloneOnly, pull, reclone
	Branch   string `yaml:"branch"`
	Tag      string `yaml:"tag"`
	OnlyIf   string `yaml:"onlyIf"`
	NotIf    string `yaml:"notIf"`
	ExitIf   string `yaml:"exitIf"`
}

func (g *Git) Setup() {
	g.Repo = RenderEnvString(g.Repo)
	g.Location = RenderEnvString(g.Location)
	g.Branch = RenderEnvString(g.Branch)
	g.Tag = RenderEnvString(g.Tag)
	g.OnlyIf = RenderEnvString(g.OnlyIf)
	g.NotIf = RenderEnvString(g.NotIf)
	g.ExitIf = RenderEnvString(g.ExitIf)
	if g.Mode == "" {
		g.Mode = "pull"
	}
	// make the destination directory without the last path element
	target := path.Dir(g.Location)
	err := os.MkdirAll(target, 0755)
	if err != nil {
		log.Error().Err(err).Msg("failed to create git destination directory for git clone")
	}
}

// Execute runs the command.
func (g *Git) Execute() error {
	g.Setup()
	/* We do not replace command envars like the other functions, this is intended to be a raw command */
	if !CanContinue(g.OnlyIf, g.NotIf, g.ExitIf, "") {
		return nil
	}
	log.Info().Msgf("git: %s to %s", g.Repo, g.Location)
	// if directory exists and it contains a .git directory, just return
	if _, err := os.Stat(path.Join(g.Location, ".git")); err == nil {
		if g.Mode == "reclone" {
			// since the mode is reclone we need to remove the directory and reclone
			err = os.RemoveAll(g.Location)
			if err != nil {
				log.Error().Err(err).Msg("failed to remove git repo for reclone")
				return err
			}
		}
		if g.Mode == "pull" {
			// since the mode is pull we need to pull the repo
			r, err := git.PlainOpen(g.Location)
			if err != nil {
				log.Error().Err(err).Msg("failed to open git repo for pull")
				return err
			}
			w, err := r.Worktree()
			if err != nil {
				log.Error().Err(err).Msg("failed to get worktree for git repo")
				return err
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
				log.Error().Err(err).Msg("failed to pull git repo")
				return err
			}
			log.Info().Msgf("git repo pulled: %s", g.Location)
			return nil
		}
		log.Info().Msgf("git repo already exists: %s", g.Location)
		return nil
	}

	r, err := git.PlainClone(g.Location, false, &git.CloneOptions{
		URL:          g.Repo,
		Progress:     os.Stdout,
		SingleBranch: g.Branch != "",
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to clone repo")
		return err
	}
	log.Info().Msgf("git cloned: %s to %s", g.Repo, g.Location)
	if g.Branch != "" {
		headRef, err := r.Head()
		if err != nil {
			log.Error().Err(err).Msg("failed to get HEAD reference")
			return err
		}

		ref := plumbing.NewHashReference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s,"+g.Branch)), headRef.Hash())
		err = r.Storer.SetReference(ref)
		if err != nil {
			log.Error().Err(err).Msg("failed to set branch reference")
			return err
		}
	}
	if g.Tag != "" {
		ref, err := r.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", g.Tag)), true)
		if err != nil {
			log.Error().Err(err).Msg("failed to get tag reference")
			return err
		}
		err = r.Storer.SetReference(ref)
		if err != nil {
			log.Error().Err(err).Msg("failed to set tag reference")
			return err
		}
	}

	return nil
}
