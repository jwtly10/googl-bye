package search

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/mock"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
	"github.com/jwtly10/googl-bye/internal/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestFindRepositories(t *testing.T) {
	// Create a mock client
	mockClient := &mock.MockGithubClient{
		MockSearchRepositories: func(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error) {
			return []*github.Repository{
					{
						Name: github.String("repo1"),
						Owner: &github.User{
							Login: github.String("owner1"),
						},
						URL:             github.String("https://api.github.com/repos/owner1/repo1"),
						Language:        github.String("Golang"),
						StargazersCount: github.Int(30),
						ForksCount:      github.Int(0),
						PushedAt:        &github.Timestamp{Time: time.Now()},
					},
					{
						Name: github.String("repo2"),
						Owner: &github.User{
							Name:  github.String("Owner Two"),
							Login: github.String("owner2"),
						},
						URL:             github.String("https://api.github.com/repos/owner2/repo2"),
						Language:        github.String("Java"),
						StargazersCount: github.Int(3921),
						ForksCount:      github.Int(901),
						PushedAt:        &github.Timestamp{Time: time.Now()},
					},
				},
				&github.Response{NextPage: 0},
				nil
		},
		MockCheckRateLimit: func(ctx context.Context) (*github.RateLimits, error) {
			return &github.RateLimits{
				Core: &github.Rate{
					Limit:     10,
					Remaining: 10,
					Reset: github.Timestamp{
						Time: time.Now(),
					},
				},
				Search: &github.Rate{
					Limit:     10,
					Remaining: 10,
					Reset: github.Timestamp{
						Time: time.Now(),
					},
				},
			}, nil
		},
	}

	container, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())
	log := common.NewLogger(false, zapcore.DebugLevel)

	searchRepo := repository.NewSearchParamRepository(db)
	repoRepo := repository.NewRepoRepository(db)

	// Create a GithubSearch instance with the mock client
	gs := &GithubSearch{
		client:     mockClient,
		config:     &common.Config{},
		log:        log,
		searchRepo: searchRepo,
		repoRepo:   repoRepo,
	}

	// Call FindRepositories
	query := &models.SearchParamsModel{
		Name:           "unit_test_1",
		StartPage:      0,
		CurrentPage:    0,
		PagesToProcess: 1,
	}
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
	assert.Equal(t, "Golang", repos[0].Language)
	assert.Equal(t, 30, repos[0].Stars)
	assert.Equal(t, 0, repos[0].Forks)
	assert.NotNil(t, repos[0].LastPush)

	// Assert correct parsing of second repo
	assert.Equal(t, "repo2", repos[1].Name)
	assert.Equal(t, "owner2", repos[1].Author)
	assert.Equal(t, "https://github.com/owner2/repo2", repos[1].GhUrl)
	assert.Equal(t, "https://github.com/owner2/repo2.git", repos[1].CloneUrl)
	assert.Equal(t, "https://api.github.com/repos/owner2/repo2", repos[1].ApiUrl)
	assert.Equal(t, "Java", repos[1].Language)
	assert.Equal(t, 3921, repos[1].Stars)
	assert.Equal(t, 901, repos[1].Forks)
	assert.NotNil(t, repos[1].LastPush)
}
