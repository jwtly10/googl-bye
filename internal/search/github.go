package search

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/errors"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
)

type GithubSearch struct {
	client     common.GithubClientI
	config     *common.Config
	log        common.Logger
	searchRepo repository.SearchParamRepository
	repoRepo   repository.RepoRepository
}

func NewGithubSearch(config *common.Config, log common.Logger, searchRepo repository.SearchParamRepository, repoRepo repository.RepoRepository) *GithubSearch {
	ghClient := common.NewGitHubClient(config.GHToken, log)

	return &GithubSearch{
		client:     ghClient,
		config:     config,
		log:        log,
		searchRepo: searchRepo,
		repoRepo:   repoRepo,
	}
}

func (ghs *GithubSearch) FindUsers(ctx context.Context, userName string) ([]*github.User, error) {
	users, _, err := ghs.client.SearchForUser(ctx, userName)
	if err != nil {
		ghs.log.Errorf("Error searching users: %v", err)
		return nil, err
	}
	ghs.log.Infof("Found %d users from API", len(users))

	return users, nil
}

func (ghs *GithubSearch) FindRepositories(ctx context.Context, params *models.SearchParamsModel) (allRepos []models.RepositoryModel, err error) {
	defer func() {
		if r := recover(); r != nil {
			ghs.log.Errorf("Panic occurred in FindRepositories: %v", r)
			err = errors.NewInternalError(fmt.Sprintf("panic occurred: %v", r))
			debug.PrintStack()
		}
	}()

	// Log rate limit before continuing
	res, err := ghs.client.CheckRateLimit(ctx)
	if err != nil {
		ghs.log.Errorf("Error checking rate limit: %v", err)
		return nil, err
	}
	ghs.log.Infof("Current search rate limits: %d/%d - Resets: %v", res.Search.Remaining, res.Search.Limit, res.Search.Reset)
	ghs.log.Infof("Current core rate limits: %d/%d - Resets: %v", res.Core.Remaining, res.Core.Limit, res.Core.Reset)

	currentPage := params.StartPage
	endPage := params.StartPage + params.PagesToProcess

	for currentPage < endPage {
		params.Opts.Page = currentPage
		ghRepos, resp, err := ghs.client.SearchRepositories(ctx, params.Query, &params.Opts)
		if err != nil {
			ghs.log.Errorf("Error searching repositories: %v", err)
			return nil, err
		}
		ghs.log.Infof("Found %d repos from API (Page %d)", len(ghRepos), currentPage)

		for _, repo := range ghRepos {
			repoName := getStringOrEmpty(repo.Name)

			parsedRepo, parseErr := ghs.parseRepo(repo)
			if parseErr != nil {
				ghs.log.Errorf("Error parsing repo %s: %v", repoName, parseErr)
				continue
			}

			allRepos = append(allRepos, parsedRepo)
		}

		if resp.NextPage == 0 {
			break // No more pages
		}
		currentPage = resp.NextPage
	}

	ghs.log.Infof("Found %d repos total", len(allRepos))
	return allRepos, nil
}

func (ghs *GithubSearch) parseRepo(repo *github.Repository) (models.RepositoryModel, error) {
	defer func() {
		if r := recover(); r != nil {
			ghs.log.Errorf("Panic occurred while parsing repo: %v: %v", repo, r)
			debug.PrintStack()
		}
	}()

	repoName := getStringOrEmpty(repo.Name)
	return models.RepositoryModel{
		Name:     repoName,
		Author:   getOwnerName(repo.Owner),
		GhUrl:    fmt.Sprintf("https://github.com/%s/%s", getStringOrEmpty(repo.Owner.Login), repoName),
		CloneUrl: fmt.Sprintf("https://github.com/%s/%s.git", getStringOrEmpty(repo.Owner.Login), repoName),
		ApiUrl:   getStringOrEmpty(repo.URL),
		Language: getStringOrEmpty(repo.Language),
		Stars:    getIntOrZero(repo.StargazersCount),
		Forks:    getIntOrZero(repo.ForksCount),
		LastPush: getTimeOrZero(&repo.PushedAt.Time),
	}, nil
}

func (ghs *GithubSearch) saveBatchedRepos(batch []models.RepositoryModel) error {
	for _, repo := range batch {
		err := ghs.repoRepo.CreateRepo(&repo)
		if err != nil {
			ghs.log.Errorf("[%s] Error saving repo: %v", fmt.Sprintf("%s/%s", repo.Author, repo.Name), err)
			return err
		}
	}

	return nil
}

func (ghs *GithubSearch) saveSearchParams(params *models.SearchParamsModel) error {
	ghs.log.Infof("Updating params for search '%s'", params.Name)
	err := ghs.searchRepo.SaveSearchParams(params)
	if err != nil {
		return err
	}
	return nil
}

func getStringOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func getIntOrZero(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

func getTimeOrZero(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
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
