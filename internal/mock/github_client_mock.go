package mock

import (
	"context"
	"github.com/google/go-github/v39/github"
)

type MockGithubClient struct {
	MockSearchRepositories func(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, error)
}

func (m *MockGithubClient) SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, error) {
	return m.MockSearchRepositories(ctx, query, opts)
}
