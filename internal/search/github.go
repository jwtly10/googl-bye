package search

import (
	"context"
	"fmt"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googlbye/internal/common"
	"github.com/jwtly10/googlbye/internal/models"
)

type GithubClientI interface {
	SearchRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, error)
}

type GithubSearch struct {
	client GithubClientI
	config *common.Config
	log    common.Logger
}

func NewGithubSearch(config *common.Config, log common.Logger) *GithubSearch {
	ghClient := common.NewGitHubClient(config.GHToken)

	return &GithubSearch{
		client: ghClient,
		config: config,
		log:    log,
	}
}

func (ghs *GithubSearch) FindRepositories(ctx context.Context, query string, opts *github.SearchOptions) ([]models.RepoModel, error) {
	ghRepos, err := ghs.client.SearchRepositories(ctx, query, opts)
	if err != nil {
		ghs.log.Errorf("Error searching repositories: %v", err)
		return nil, nil
	}

	ghs.log.Debugf("Repos from API: %v", ghRepos)

	var repos []models.RepoModel
	for _, repo := range ghRepos {
		repoName := getStringOrEmpty(repo.Name)

		ghs.log.Debugf("Trying to parse: %v", repo)
		parsedRepo := models.RepoModel{
			Name:     repoName,
			Author:   getOwnerName(repo.Owner),
			GhUrl:    fmt.Sprintf("https://github.com/%s/%s", *repo.Owner.Login, repoName),
			CloneUrl: fmt.Sprintf("https://github.com/%s/%s.git", *repo.Owner.Login, repoName),
			ApiUrl:   getStringOrEmpty(repo.URL),
		}

		repos = append(repos, parsedRepo)
	}

	ghs.log.Infof("Parsed %d repos", len(repos))

	for _, r := range repos {
		ghs.log.Debugf("%v", r)
	}

	return repos, nil
}

func getStringOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// Gets either the login name or the nicely formatted owners name
func getOwnerName(owner *github.User) string {
	if owner != nil && owner.Name != nil {
		return *owner.Name
	}
	if owner != nil && owner.Login != nil {
		return *owner.Login
	}
	return ""
}
