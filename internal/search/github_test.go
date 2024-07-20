package search

import (
	"context"
	"testing"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/mock"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestFindRepositories(t *testing.T) {
	// Create a mock client
	mockClient := &mock.MockGithubClient{
		MockSearchRepositories: func(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, error) {
			return []*github.Repository{
				{
					Name: github.String("repo1"),
					Owner: &github.User{
						Login: github.String("owner1"),
					},
					URL: github.String("https://api.github.com/repos/owner1/repo1"),
				},
				{
					Name: github.String("repo2"),
					Owner: &github.User{
						Name:  github.String("Owner Two"),
						Login: github.String("owner2"),
					},
					URL: github.String("https://api.github.com/repos/owner2/repo2"),
				},
			}, nil
		},
	}

	// Create a GithubSearch instance with the mock client
	gs := &GithubSearch{
		client: mockClient,
		config: &common.Config{},
		log:    common.NewLogger(false, zapcore.DebugLevel),
	}

	// Call FindRepositories
	query := &models.SearchParams{}
	repos, err := gs.FindRepositories(context.Background(), query)

	// Assert no error
	assert.NoError(t, err)

	// Assert correct number of repos
	assert.Len(t, repos, 2)

	// Assert correct parsing of first repo
	assert.Equal(t, "repo1", repos[0].Name)
	assert.Equal(t, "owner1", repos[0].Author)
	assert.Equal(t, "https://github.com/owner1/repo1", repos[0].GhUrl)
	assert.Equal(t, "https://github.com/owner1/repo1.git", repos[0].CloneUrl)
	assert.Equal(t, "https://api.github.com/repos/owner1/repo1", repos[0].ApiUrl)

	// Assert correct parsing of second repo
	assert.Equal(t, "repo2", repos[1].Name)
	assert.Equal(t, "Owner Two", repos[1].Author) // Uses Name instead of Login
	assert.Equal(t, "https://github.com/owner2/repo2", repos[1].GhUrl)
	assert.Equal(t, "https://github.com/owner2/repo2.git", repos[1].CloneUrl)
	assert.Equal(t, "https://api.github.com/repos/owner2/repo2", repos[1].ApiUrl)
}
