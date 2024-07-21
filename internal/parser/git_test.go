package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestGitCmdLineClone(t *testing.T) {

	logger := common.NewLogger(false, zapcore.DebugLevel)
	gitCmdLine := NewGitCmdLine(logger)

	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "git-clone-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Use a public repository URL for testing
	repoURL := "https://github.com/jwtly10/googl-bye-test.git"
	destination := filepath.Join(tmpDir, "repo")

	// Execute
	err = gitCmdLine.Clone(repoURL, destination)

	// Assert
	assert.NoError(t, err)

	// Check if README.md exists
	readmePath := filepath.Join(destination, "README.md")
	_, err = os.Stat(readmePath)
	assert.NoError(t, err, "README.md should exist in the cloned repository")
}
