package parser

import (
	"testing"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestParseRepository(t *testing.T) {
	logger := common.NewLogger(false, zapcore.DebugLevel)
	git := NewGitCmdLine(logger)

	parser := NewRepoParser(git, logger)

	repo := models.RepositoryModel{
		Name:     "googl-bye-test",
		Author:   "jwtly10",
		CloneUrl: "https://github.com/jwtly10/googl-bye-test.git",
	}

	links, err := parser.ParseRepository(repo)
	if err != nil {
		t.Errorf("expected no error when parsing repository but got %v", err)
	}

	assert.Len(t, links, 4)

	assert.Equal(t, "http://goo.gl/Y5VIoG", links[0].Url)
	assert.Equal(t, "http://google.com/", links[0].ExpandedUrl)
	assert.Equal(t, "README.md", links[0].File)
	assert.Equal(t, 5, links[0].LineNumber)
	assert.Equal(t, "https://github.com/jwtly10/googl-bye-test/blob/main/README.md?plain=1#L5", links[0].GithubUrl)

	assert.Equal(t, "https://goo.gl/aoDfac", links[1].Url)
	assert.Equal(t, "http://www.timesofisrael.com/in-central-asia-netanyahu-scores-diplomatic-home-run-in-irans-backyard/", links[1].ExpandedUrl)
	assert.Equal(t, "README.md", links[1].File)
	assert.Equal(t, 7, links[1].LineNumber)
	assert.Equal(t, "https://github.com/jwtly10/googl-bye-test/blob/main/README.md?plain=1#L7", links[1].GithubUrl)

	assert.Equal(t, "http://goo.gl/Y5VIoG", links[2].Url)
	assert.Equal(t, "http://google.com/", links[2].ExpandedUrl)
	assert.Equal(t, "main.go", links[2].File)
	assert.Equal(t, 4, links[2].LineNumber)
	assert.Equal(t, "https://github.com/jwtly10/googl-bye-test/blob/main/main.go#L4", links[2].GithubUrl)

	assert.Equal(t, "http://goo.gl/forms/xm5KFo35tu", links[3].Url)
	assert.Contains(t, links[3].ExpandedUrl, "https://docs.google.com/forms/d/")
	assert.Equal(t, "main.go", links[3].File)
	assert.Equal(t, 7, links[3].LineNumber)
	assert.Equal(t, "https://github.com/jwtly10/googl-bye-test/blob/main/main.go#L7", links[3].GithubUrl)
}
