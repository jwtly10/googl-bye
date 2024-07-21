package search

import (
	"context"
	"fmt"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
)

type RepoSearch struct {
	params   *models.SearchParams
	log      common.Logger
	repoRepo repository.RepoRepository
	gh       *GithubSearch
	cache    *map[string]bool
}

func NewRepoSearch(params *models.SearchParams, config *common.Config, log common.Logger, repoRepo repository.RepoRepository, cache *map[string]bool) *RepoSearch {
	gh := NewGithubSearch(config, log, cache)
	return &RepoSearch{
		params:   params,
		gh:       gh,
		log:      log,
		repoRepo: repoRepo,
	}
}

func (rs *RepoSearch) StartSearch(ctx context.Context) {
	rs.log.Infof("Running search for params 'Query: %s', 'Params: %v'", rs.params.Query, *rs.params.Opts)
	repos, err := rs.gh.FindRepositories(ctx, rs.params)
	if err != nil {
		rs.log.Errorf("Error fetching repositories: %v", err)
	}

	for _, repo := range repos {
		err := rs.repoRepo.CreateRepo(&repo)
		if err != nil {
			rs.log.Errorf("Error creating repo in db: %v", err)
		}
		// Update the cache
		(*rs.cache)[fmt.Sprintf("%s/%s", repo.Author, repo.Name)] = true
	}
}
