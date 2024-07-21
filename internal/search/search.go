package search

import (
	"context"
	"fmt"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
)

type RepoSearch struct {
	params   *models.SearchParamsModel
	log      common.Logger
	gh       *GithubSearch
	cache    *map[string]bool
	repoRepo repository.RepoRepository
}

func NewRepoSearch(params *models.SearchParamsModel, config *common.Config, log common.Logger, repoRepo repository.RepoRepository, cache *map[string]bool, searchRepo repository.SearchParamRepository) *RepoSearch {
	gh := NewGithubSearch(config, log, cache, searchRepo)
	return &RepoSearch{
		params:   params,
		gh:       gh,
		log:      log,
		repoRepo: repoRepo,
		cache:    cache,
	}
}

func (rs *RepoSearch) StartSearch(ctx context.Context) {
	rs.log.Infof("Running search for params 'Query: %s', 'Params: %v' 'StartPage': %d, 'CurrentPage': %d, 'PagesToProcess': %d", rs.params.Query, rs.params.Opts, rs.params.StartPage, rs.params.CurrentPage, rs.params.PagesToProcess)
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
		rs.log.Infof("[%s] Repo saved to DB", fmt.Sprintf("%s/%s", repo.Author, repo.Name))
	}
}
