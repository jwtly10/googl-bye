package parser

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jwtly10/googl-bye/internal/common"
)

type GitCmdLineI interface {
	Clone(url, destination string) (string, error)
}

type GitCmdLine struct {
	log common.Logger
}

func NewGitCmdLine(log common.Logger) *GitCmdLine {
	return &GitCmdLine{
		log: log,
	}
}

func (g *GitCmdLine) Clone(url, destination string) (string, error) {
	// Clone the repository
	g.log.Infof("Cloning repo '%s' into '%s'", url, destination)
	cloneCmd := exec.Command("git", "clone", "--depth", "1", url, destination)
	output, err := cloneCmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("failed to clone repository: %w\nOutput: %s", err, string(output))
    }

	// Get the current branch
	branchCmd := exec.Command("git", "-C", destination, "rev-parse", "--abbrev-ref", "HEAD")
	output, err = branchCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w\nOutput: %s", err, string(output))
	}

	// Trim any whitespace from the branch name
	branch := strings.TrimSpace(string(output))

	g.log.Infof("Cloned repo '%s' into '%s'. Current branch: %s", url, destination, branch)

	return branch, nil
}
