package common

import (
	"context"
	"fmt"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type GithubClientI interface {
	SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error)
	SearchForUser(ctx context.Context, username string) ([]*github.User, *github.Response, error)
	CheckRateLimit(ctx context.Context) (*github.RateLimits, error)
}

type GithubClient struct {
	client *github.Client
}

func NewGitHubClient(token string) *GithubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GithubClient{client: client}
}

func (gc *GithubClient) SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error) {
	result, response, err := gc.client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("error searching repositories: %v", err)
	}

	return result.Repositories, response, nil
}

func (gc *GithubClient) SearchForUser(ctx context.Context, username string) ([]*github.User, *github.Response, error) {
	result, response, err := gc.client.Search.Users(ctx, username, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error searching users: %v", err)
	}

	return result.Users, response, nil
}

func (gc *GithubClient) CheckRateLimit(ctx context.Context) (*github.RateLimits, error) {
	result, _, err := gc.client.RateLimits(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking rate limit: %v", err)
	}

	return result, nil
}
