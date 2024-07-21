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

	searchRepo := repository.NewSearchParamRepository(db)
	repoRepo := repository.NewRepoRepository(db)

	// Create a GithubSearch instance with the mock client
	cache := make(map[string]bool)
	gs := &GithubSearch{
		client:     mockClient,
		config:     &common.Config{},
		log:        common.NewLogger(false, zapcore.DebugLevel),
		repoCache:  &cache,
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

	// Assert correct parsing of second repo
	assert.Equal(t, "repo2", repos[1].Name)
	assert.Equal(t, "owner2", repos[1].Author)
	assert.Equal(t, "https://github.com/owner2/repo2", repos[1].GhUrl)
	assert.Equal(t, "https://github.com/owner2/repo2.git", repos[1].CloneUrl)
	assert.Equal(t, "https://api.github.com/repos/owner2/repo2", repos[1].ApiUrl)
}

func TestFindRepositoriesCacheHit(t *testing.T) {
	// Create a mock client
	mockClient := &mock.MockGithubClient{
		MockSearchRepositories: func(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error) {
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

	searchRepo := repository.NewSearchParamRepository(db)
	repoRepo := repository.NewRepoRepository(db)

	// Create a GithubSearch instance with the mock client
	cache := make(map[string]bool)
	// Adding one of the repos to cache
	cache["owner1/repo1"] = true
	gs := &GithubSearch{
		client:     mockClient,
		config:     &common.Config{},
		log:        common.NewLogger(false, zapcore.DebugLevel),
		repoCache:  &cache,
		searchRepo: searchRepo,
		repoRepo:   repoRepo,
	}

	// Call FindRepositories
	query := &models.SearchParamsModel{
		Name:  "unit_test_2",
		Query: "stars:>10000",
		Opts: github.SearchOptions{
			Sort:  "stars",
			Order: "desc",
			ListOptions: github.ListOptions{
				PerPage: 2,
				Page:    0,
			},
		},
		PagesToProcess: 1,
		CurrentPage:    1,
		StartPage:      1,
	}

	repos, err := gs.FindRepositories(context.Background(), query)

	// Assert no error
	assert.NoError(t, err)

	// Assert correct number of repos
	// Note only 1 this time as there was a cache hit.
	assert.Len(t, repos, 1)

	// Assert correct parsing of second repo
	assert.Equal(t, "repo2", repos[0].Name)
	assert.Equal(t, "owner2", repos[0].Author)
	assert.Equal(t, "https://github.com/owner2/repo2", repos[0].GhUrl)
	assert.Equal(t, "https://github.com/owner2/repo2.git", repos[0].CloneUrl)
	assert.Equal(t, "https://api.github.com/repos/owner2/repo2", repos[0].ApiUrl)
}
