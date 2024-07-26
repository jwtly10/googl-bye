package common

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type GithubClientI interface {
	SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error)
	SearchForUser(ctx context.Context, username string) ([]*github.User, *github.Response, error)
	CheckRateLimit(ctx context.Context) (*github.RateLimits, error)
	CreateIssue(ctx context.Context, owner, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}

type GithubClient struct {
	client *github.Client
	log    Logger
}

func NewGitHubClient(token string, log Logger) *GithubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &GithubClient{client: client, log: log}
}

func (gc *GithubClient) CreateIssue(ctx context.Context, owner, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	result, response, err := gc.client.Issues.Create(ctx, owner, repo, issue)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating issue: %v", err)
	}

	return result, response, nil
}

func (gc *GithubClient) SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, *github.Response, error) {
	result, response, err := gc.client.Search.Repositories(ctx, query, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("error searching repositories: %v", err)
	}

	return result.Repositories, response, nil
}

func (gc *GithubClient) SearchForUser(ctx context.Context, username string) ([]*github.User, *github.Response, error) {
	result, response, err := gc.client.Search.Users(ctx, username, &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error searching users: %v", err)
	}

	var wg sync.WaitGroup
	fullUsers := make([]*github.User, len(result.Users))
	errors := make([]error, len(result.Users))

	for i, user := range result.Users {
		wg.Add(1)
		go func(i int, user *github.User) {
			defer wg.Done()
			fullUser, _, err := gc.client.Users.Get(ctx, user.GetLogin())
			if err != nil {
				gc.log.Warnf("error fetching details for user %s: %v", user.GetLogin(), err)
				errors[i] = err
				return
			}
			fullUsers[i] = fullUser
		}(i, user)
	}

	wg.Wait()

	validUsers := make([]*github.User, 0, len(fullUsers))
	for i, user := range fullUsers {
		if user != nil {
			validUsers = append(validUsers, user)
		} else if errors[i] != nil {
			gc.log.Warnf("error fetching details for user: %v", errors[i])
		}
	}

	rateLimit, err := gc.CheckRateLimit(ctx)
	if err != nil {
		gc.log.Errorf("Error checking rate limit: %v", err)
	} else {
		gc.log.Infof("Current search rate limits: %d/%d - Resets: %v", rateLimit.Search.Remaining, rateLimit.Search.Limit, rateLimit.Search.Reset)
		gc.log.Infof("Current core rate limits: %d/%d - Resets: %v", rateLimit.Core.Remaining, rateLimit.Core.Limit, rateLimit.Core.Reset)
	}

	return validUsers, response, nil
}

func (gc *GithubClient) CheckRateLimit(ctx context.Context) (*github.RateLimits, error) {
	result, _, err := gc.client.RateLimits(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking rate limit: %v", err)
	}

	return result, nil
}
