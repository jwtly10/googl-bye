package search

import (
	"context"
	"fmt"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
)

type GithubSearch struct {
	client     common.GithubClientI
	config     *common.Config
	log        common.Logger
	repoCache  *map[string]bool
	searchRepo repository.SearchParamRepository
}

func NewGithubSearch(config *common.Config, log common.Logger, repoCache *map[string]bool) *GithubSearch {
	ghClient := common.NewGitHubClient(config.GHToken)

	return &GithubSearch{
		client:    ghClient,
		config:    config,
		log:       log,
		repoCache: repoCache,
	}
}

func (ghs *GithubSearch) FindRepositories(ctx context.Context, params *models.SearchParamsModel) ([]models.RepositoryModel, error) {
	// Log rate limit before continuing
	res, err := ghs.client.CheckRateLimit(ctx)
	if err != nil {
		ghs.log.Errorf("Error checking rate limit: %v", err)
		return nil, err
	}

	ghs.log.Infof("Current search rate limits: %d/%d - Resets: %v", res.Search.Remaining, res.Search.Limit, res.Search.Reset)
	ghs.log.Infof("Current core rate limits: %d/%d - Resets: %v", res.Core.Remaining, res.Core.Limit, res.Core.Reset)

	var allRepos []models.RepositoryModel
	var cacheHits int
	currentPage := params.StartPage
	endPage := params.StartPage + params.PagesToProcess

	for currentPage < endPage {
		params.Opts.Page = currentPage

		ghRepos, resp, err := ghs.client.SearchRepositories(ctx, params.Query, &params.Opts)
		if err != nil {
			ghs.log.Errorf("Error searching repositories: %v", err)
			return nil, err
		}

		ghs.log.Debugf("Repos from API (Page %d): %v", currentPage, ghRepos)
		ghs.log.Infof("Found %d repos from API (Page %d)", len(ghRepos), currentPage)

		for _, repo := range ghRepos {
			repoName := getStringOrEmpty(repo.Name)

			// check the cache
			if (*ghs.repoCache)[fmt.Sprintf("%s/%s", *repo.Owner.Login, repoName)] {
				// If this repo is already in cache, skip doing anything with it
				ghs.log.Infof("[%s] Cache hit", fmt.Sprintf("%s/%s", *repo.Owner.Login, repoName))
				cacheHits++
				continue
			}

			ghs.log.Debugf("Trying to parse data from repo : %v", repo)
			parsedRepo := models.RepositoryModel{
				Name:     repoName,
				Author:   getOwnerName(repo.Owner),
				GhUrl:    fmt.Sprintf("https://github.com/%s/%s", *repo.Owner.Login, repoName),
				CloneUrl: fmt.Sprintf("https://github.com/%s/%s.git", *repo.Owner.Login, repoName),
				ApiUrl:   getStringOrEmpty(repo.URL),
			}
			allRepos = append(allRepos, parsedRepo)
		}

		// Update the search params with the current page
		params.CurrentPage = currentPage
		err = ghs.saveSearchParams(params)
		if err != nil {
			ghs.log.Errorf("Error saving search params: %v", err)
			continue
		}

		if resp.NextPage == 0 {
			break // No more pages
		}
		currentPage = resp.NextPage
	}

	ghs.log.Infof("Found %d repos to save (%d cache hits)", len(allRepos), cacheHits)
	for _, r := range allRepos {
		ghs.log.Debugf("%v", r)
	}
	return allRepos, nil
}

func (ghs *GithubSearch) saveSearchParams(params *models.SearchParamsModel) error {
	return nil
}

// // TODO implement pagination, and keep track of the last page we were on for a given search
// func (ghs *GithubSearch) FindRepositories(ctx context.Context, params *models.SearchParams) ([]models.RepositoryModel, error) {

// 	// Log rate limit before continuing
// 	res, err := ghs.client.CheckRateLimit(ctx)
// 	if err != nil {
// 		ghs.log.Errorf("Error checking rate limit: %v", err)
// 		return nil, err
// 	}

// 	ghs.log.Infof("Current search rate limits: %d/%d - Resets: %v", res.Search.Remaining, res.Search.Limit, res.Search.Reset)
// 	ghs.log.Infof("Current core rate limits: %d/%d - Resets: %v", res.Core.Remaining, res.Core.Limit, res.Core.Reset)

// 	ghRepos, err := ghs.client.SearchRepositories(ctx, params.Query, params.Opts)
// 	if err != nil {
// 		ghs.log.Errorf("Error searching repositories: %v", err)
// 		return nil, err
// 	}

// 	ghs.log.Debugf("Repos from API: %v", ghRepos)

// 	var repos []models.RepositoryModel
// 	for _, repo := range ghRepos {
// 		repoName := getStringOrEmpty(repo.Name)

// 		// check the cache
// 		if (*ghs.repoCache)[fmt.Sprintf("%s/%s", *repo.Owner.Login, repoName)] {
// 			// If this repo is already in cache, skip doing anything with it
// 			ghs.log.Infof("[%s] Cache hit", fmt.Sprintf("%s/%s", *repo.Owner.Login, repoName))
// 			continue
// 		}

// 		ghs.log.Debugf("Trying to parse: %v", repo)
// 		parsedRepo := models.RepositoryModel{
// 			Name:     repoName,
// 			Author:   getOwnerName(repo.Owner),
// 			GhUrl:    fmt.Sprintf("https://github.com/%s/%s", *repo.Owner.Login, repoName),
// 			CloneUrl: fmt.Sprintf("https://github.com/%s/%s.git", *repo.Owner.Login, repoName),
// 			ApiUrl:   getStringOrEmpty(repo.URL),
// 		}

// 		repos = append(repos, parsedRepo)
// 	}

// 	ghs.log.Infof("Saved %d repos to DB", len(repos))

// 	for _, r := range repos {
// 		ghs.log.Debugf("%v", r)
// 	}

// 	return repos, nil
// }

func getStringOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// Gets either the login name or the nicely formatted owners name
// TODO: Removing the nice name for now, only returning the login
func getOwnerName(owner *github.User) string {
	// if owner != nil && owner.Name != nil {
	// 	return *owner.Name
	// }
	if owner != nil && owner.Login != nil {
		return *owner.Login
	}
	return ""
}
