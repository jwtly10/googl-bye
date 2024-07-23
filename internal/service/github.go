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
	"github.com/jwtly10/googl-bye/internal/search"
	"github.com/jwtly10/googl-bye/internal/utils"
)

type GithubService struct {
	log common.Logger
	ghs search.GithubSearch
}

func NewGithubService(ghc search.GithubSearch, log common.Logger) *GithubService {
	return &GithubService{
		log: log,
		ghs: ghc,
	}
}

func (gs *GithubService) GithubSearch(r *http.Request) ([]models.RepositoryModel, error) {
	searchParams, err := gs.validateBodyFromRequest(r)
	if err != nil {
		return nil, err
	}

	gs.log.Infof("%v", searchParams)

	// Force max repo size of 50MB (TODO: Review)
	if !strings.Contains(searchParams.Query, "size:<=50000") {
		searchParams.Query = searchParams.Query + " size:<=50000"
	}

	githubOpts := &github.SearchOptions{
		Sort:  searchParams.Opts.Sort,
		Order: searchParams.Opts.Order,
		ListOptions: github.ListOptions{
			Page:    searchParams.StartPage,
			PerPage: 100,
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

	return res, nil
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

	gs.log.Debugf("Decoded struct: %+v\n", searchParams)

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
