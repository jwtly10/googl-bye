package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/errors"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
	"github.com/jwtly10/googl-bye/internal/search"
	"github.com/jwtly10/googl-bye/internal/utils"
)

type GithubService struct {
	log          common.Logger
	ghs          search.GithubSearch
	repoLinkRepo repository.RepoLinkRepository
}

func NewGithubService(ghc search.GithubSearch, r repository.RepoLinkRepository, log common.Logger) *GithubService {
	return &GithubService{
		log:          log,
		ghs:          ghc,
		repoLinkRepo: r,
	}
}

func (gs *GithubService) GithubCreateIssue(r *http.Request) (string, error) {
	repoId, err := gs.getRepoIdFromBody(r)
	if err != nil {
		return "", err
	}

	gs.log.Infof("Creating issue for repo: %s", repoId)

	repoDetails, err := gs.repoLinkRepo.GetRepoLinksById(repoId)
	if err != nil {
		return "", errors.NewInternalError(fmt.Sprintf("error getting repo from ID: %v", err))
	}

	issue, err := gs.ghs.CreateIssueFromRepo(r.Context(), repoDetails)
	if err != nil {
		return "", errors.NewInternalError(fmt.Sprintf("error creating issue: %v", err))
	}

	gs.log.Infof("Issue created successfully: %v", issue.GetURL())

	if issue != nil {
		return issue.GetHTMLURL(), nil
	} else {
		return "", errors.NewInternalError("error creating issue. Issue struct is nil")
	}
}

func (gs *GithubService) GithubSearchRepos(r *http.Request) ([]models.RepositoryModel, error) {
	searchParams, err := gs.validateBodyFromRequest(r)
	if err != nil {
		return nil, err
	}

	gs.log.Infof("Github search Params: %v", searchParams)

	// Force max repo size of 50MB (TODO: Review)
	if !strings.Contains(searchParams.Query, "size:<=50000") {
		searchParams.Query = searchParams.Query + " size:<=50000"
	}
	// Force minimum stars > 300 (TODO: Review)
	if !strings.Contains(searchParams.Query, "stars") {
		searchParams.Query = searchParams.Query + " stars:>300"
	}

	githubOpts := &github.SearchOptions{
		Sort:  searchParams.Opts.Sort,
		Order: searchParams.Opts.Order,
		ListOptions: github.ListOptions{
			Page:    searchParams.StartPage,
			PerPage: 50,
		},
	}

	searchParams.Opts = *githubOpts
	// TODO: Remove Name from table
	searchParams.Name = utils.GenerateShortID()

	gs.log.Infof("Running search for params 'Query: %s', 'Params: %v' 'StartPage': %d, 'CurrentPage': %d, 'PagesToProcess': %d", searchParams.Query, searchParams.Opts, searchParams.StartPage, searchParams.CurrentPage, searchParams.PagesToProcess)

	res, err := gs.ghs.FindRepositories(r.Context(), searchParams)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Sprintf("error finding repositories: %v", err.Error()))
	}

	if len(res) == 0 {
		return []models.RepositoryModel{}, nil
	}

	return res, nil
}

func (gs *GithubService) GithubSearchReposForUser(r *http.Request) ([]models.RepositoryModel, error) {
	username := r.URL.Query().Get("username")
	if username == "" {
		return nil, errors.NewBadRequestError("missing required field: username")
	}

	gs.log.Infof("Github repo search for user : %v", username)

	searchParams := &models.SearchParamsModel{
		Query: fmt.Sprintf("user:%s size:<=50000", username), // Also limit to 50MB
		Opts: github.SearchOptions{
			Sort:  "updated",
			Order: "desc",
			ListOptions: github.ListOptions{
				PerPage: 50,
			},
		},
		StartPage:      0,
		CurrentPage:    0,
		PagesToProcess: 2, // At 50 per page, we only load at most 100 repos TODO: This breaks for users with more than 100 repos
	}

	gs.log.Infof("Running user repo search for params 'Query: %s', 'Params: %v' 'StartPage': %d, 'CurrentPage': %d, 'PagesToProcess': %d", searchParams.Query, searchParams.Opts, searchParams.StartPage, searchParams.CurrentPage, searchParams.PagesToProcess)

	res, err := gs.ghs.FindRepositories(r.Context(), searchParams)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Sprintf("error finding repositories: %v", err.Error()))
	}

	if len(res) == 0 {
		return []models.RepositoryModel{}, nil
	}

	return res, nil
}

func (gs *GithubService) GithubSearchUsers(r *http.Request) ([]models.GithubUser, error) {

	username := r.URL.Query().Get("username")
	if username == "" {
		return nil, errors.NewBadRequestError("missing required field: username")
	}

	gs.log.Infof("Github user search for: %v", username)

	res, err := gs.ghs.FindUsers(r.Context(), username)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Sprintf("error finding user: %v", err.Error()))
	}

	var users []models.GithubUser

	for _, user := range res {
		users = append(users, models.GithubUser{
			Id:        user.GetID(),
			Login:     user.GetLogin(),
			Url:       user.GetURL(),
			AvatarUrl: user.GetAvatarURL(),
			Name:      user.GetName(),
		})
	}

	if len(users) == 0 {
		return []models.GithubUser{}, nil
	}

	return users, nil
}

func (gs *GithubService) validateBodyFromRequest(r *http.Request) (*models.SearchParamsModel, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Sprintf("error reading request body: %v", err.Error()))
	}

	gs.log.Debugf("Raw JSON from request: %s", string(body))

	var searchParams models.SearchParamsModel
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&searchParams); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, errors.NewBadRequestError(fmt.Sprintf("invalid type for field: %v", err.Error()))
		}
		if strings.HasPrefix(err.Error(), "json: unknown field") {
			return nil, errors.NewBadRequestError(fmt.Sprintf("unknown field in request body: %v", err.Error()))
		}
		return nil, errors.NewBadRequestError(fmt.Sprintf("error decoding request body: %v", err.Error()))
	}

	gs.log.Debugf("Decoded struct: %+v", searchParams)

	// Check if all required fields are present
	missingFields := []string{}
	// See above TODO: Remove name from this model
	// if searchParams.Name == "" {
	// 	missingFields = append(missingFields, "name")
	// }
	if searchParams.Query == "" {
		missingFields = append(missingFields, "query")
	}
	if searchParams.Opts == (github.SearchOptions{}) {
		missingFields = append(missingFields, "opts")
	}
	// if searchParams.StartPage == 0 {
	// 	missingFields = append(missingFields, "startPage")
	// }
	// Dont need to validate this, its handled internally
	// if searchParams.CurrentPage == 0 {
	// 	missingFields = append(missingFields, "currentPage")
	// }
	if searchParams.PagesToProcess == 0 {
		missingFields = append(missingFields, "pagesToProcess")
	}

	if len(missingFields) > 0 {
		return nil, errors.NewBadRequestError(fmt.Sprintf("missing required fields: %s", strings.Join(missingFields, ", ")))
	}

	return &searchParams, nil
}

func (gs *GithubService) getRepoIdFromBody(r *http.Request) (int, error) {

	// The body is just {"repoId": 1} // or some other number. validate this and get the number

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, errors.NewInternalError(fmt.Sprintf("error reading request body: %v", err.Error()))
	}

	gs.log.Debugf("Raw JSON from request: %s", string(body))

	var repoId struct {
		RepoId int `json:"repoId"`
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	// We dont care about additional fields. We wont read them. We just want repo id from body

	if err := dec.Decode(&repoId); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return 0, errors.NewBadRequestError(fmt.Sprintf("invalid type for field: %v", err.Error()))
		}

		return 0, errors.NewBadRequestError(fmt.Sprintf("error decoding request body: %v", err.Error()))
	}

	gs.log.Debugf("Decoded struct: %+v", repoId)

	if repoId.RepoId == 0 {
		return 0, errors.NewBadRequestError("repoId is required")
	}

	return repoId.RepoId, nil
}
