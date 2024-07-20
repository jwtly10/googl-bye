package search

import (
	"context"

	"github.com/jwtly10/googlbye/internal/common"
	"github.com/jwtly10/googlbye/internal/models"
)

type RepoSearch struct {
	params *models.SearchParams
	log    common.Logger
	gh     *GithubSearch
}

func NewRepoSearch(params *models.SearchParams, config *common.Config, log common.Logger) *RepoSearch {
	gh := NewGithubSearch(config, log)
	return &RepoSearch{
		params: params,
		gh:     gh,
		log:    log,
	}
}

func (rs *RepoSearch) StartSearch(ctx context.Context) {
	rs.log.Infof("Running search for params %v", rs.params)
	rs.gh.FindRepositories(ctx, rs.params)
}
