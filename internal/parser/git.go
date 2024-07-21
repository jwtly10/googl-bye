package parser

import (
	"os/exec"

	"github.com/jwtly10/googl-bye/internal/common"
)

type GitCmdLineI interface {
	Clone(url, destination string) error
}

type GitCmdLine struct {
	log common.Logger
}

func NewGitCmdLine(log common.Logger) *GitCmdLine {
	return &GitCmdLine{
		log: log,
	}
}

func (g *GitCmdLine) Clone(url, destination string) error {
	// Only clone the main/most recent code. We don't need all the history
	g.log.Infof("Cloning repo '%s' into '%s'", url, destination)
	cmd := exec.Command("git", "clone", "--depth", "1", url, destination)
	return cmd.Run()
}
