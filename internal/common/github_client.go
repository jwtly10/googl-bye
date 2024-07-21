package common

import (
	"context"
	"fmt"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type GithubClientI interface {
	SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error)
	CheckRateLimit(ctx context.Context) (*github.RateLimits, error)
}

type GitHubClient struct {
	client *github.Client
}

func NewGitHubClient(token string) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GitHubClient{client: client}
}

func (gc *GitHubClient) SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error) {
	result, response, err := gc.client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("error searching repositories: %v", err)
	}

	return result.Repositories, response, nil
}

func (gc *GitHubClient) CheckRateLimit(ctx context.Context) (*github.RateLimits, error) {
	result, _, err := gc.client.RateLimits(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking rate limit: %v", err)
	}

	return result, nil
}
