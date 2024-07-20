package search

import (
	"context"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
)

type RepoSearch struct {
	params   *models.SearchParams
	log      common.Logger
	repoRepo repository.RepoRepository
	gh       *GithubSearch
}

func NewRepoSearch(params *models.SearchParams, config *common.Config, log common.Logger, repoRepo repository.RepoRepository) *RepoSearch {
	gh := NewGithubSearch(config, log)
	return &RepoSearch{
		params:   params,
		gh:       gh,
		log:      log,
		repoRepo: repoRepo,
	}
}

func (rs *RepoSearch) StartSearch(ctx context.Context) {
	rs.log.Infof("Running search for params %v", rs.params)
	repos, err := rs.gh.FindRepositories(ctx, rs.params)
	if err != nil {
		rs.log.Errorf("Error fetching repositories: %v", err)
	}

	for _, repo := range repos {
		err := rs.repoRepo.CreateRepo(&repo)
		if err != nil {
			rs.log.Errorf("Error creating repo in db: %v", err)
		}
	}
}
