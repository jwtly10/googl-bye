package mock

import (
	"context"
	"github.com/google/go-github/v39/github"
)

type MockGithubClient struct {
	MockSearchRepositories func(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error)
	MockCheckRateLimit     func(ctx context.Context) (*github.RateLimits, error)
	MockSearchForUser      func(ctx context.Context, username string) ([]*github.User, *github.Response, error)
}

func (m *MockGithubClient) SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error) {
	return m.MockSearchRepositories(ctx, query, opts)
}

func (m *MockGithubClient) CheckRateLimit(ctx context.Context) (*github.RateLimits, error) {
	return m.MockCheckRateLimit(ctx)
}

func (m *MockGithubClient) SearchForUser(ctx context.Context, username string) ([]*github.User, *github.Response, error) {
	return m.MockSearchForUser(ctx, username)
}
