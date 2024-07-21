package mock

import (
	"context"
	"github.com/google/go-github/v39/github"
)

type MockGithubClient struct {
	MockSearchRepositories func(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, error)
	MockCheckRateLimit     func(ctx context.Context) (*github.RateLimits, error)
}

func (m *MockGithubClient) SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, error) {
	return m.MockSearchRepositories(ctx, query, opts)
}

func (m *MockGithubClient) CheckRateLimit(ctx context.Context) (*github.RateLimits, error) {
	return m.MockCheckRateLimit(ctx)
}
