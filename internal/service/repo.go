package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/errors"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
)

type RepoService struct {
	log   common.Logger
	r     repository.RepoRepository
	cache *common.RepoCache
}

func NewRepoService(r repository.RepoRepository, l common.Logger, c *common.RepoCache) *RepoService {
	return &RepoService{
		r:     r,
		log:   l,
		cache: c,
	}
}

func (rs *RepoService) SaveRepo(r *http.Request) error {
	reposFromReq, err := rs.validateBodyFromRequest(r)
	if err != nil {
		return err
	}

	var cacheHits int
	var reposToSave []*models.RepositoryModel

	// For each repo. Check if cache exists for it
	// This means that the rows are already in the DB and we can just simulate that it 'saved'.
	for _, repo := range reposFromReq {
		key := fmt.Sprintf("%s/%s", repo.Author, repo.Name)
		if _, ok := rs.cache.Get(key); ok {
			rs.log.Debugf("[%s] repo found in cache. Ignoring.", key)
			cacheHits++
		} else {
			reposToSave = append(reposToSave, repo)
		}
	}

	err = rs.r.CreateRepos(reposToSave)
	if err != nil {
		return errors.NewInternalError(fmt.Sprintf("error when batch saving repositories: %v", err.Error()))
	}

	// Update the cache once all rows saved
	for _, savedRepo := range reposToSave {
		key := fmt.Sprintf("%s/%s", savedRepo.Author, savedRepo.Name)
		rs.cache.Set(key, true)
	}

	rs.log.Infof("%d repos in save request: %d saved (%d were cache hits)", len(reposFromReq), len(reposFromReq)-cacheHits, cacheHits)

	return nil
}

func (rs *RepoService) validateBodyFromRequest(r *http.Request) ([]*models.RepositoryModel, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Sprintf("error reading request body: %v", err.Error()))
	}
	rs.log.Debugf("Raw JSON from request: %s", string(body))

	var repositories []*models.RepositoryModel
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&repositories); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, errors.NewBadRequestError(fmt.Sprintf("invalid type for field: %v", err.Error()))
		}
		if strings.HasPrefix(err.Error(), "json: unknown field") {
			return nil, errors.NewBadRequestError(fmt.Sprintf("unknown field in request body: %v", err.Error()))
		}
		return nil, errors.NewBadRequestError(fmt.Sprintf("error decoding request body: %v", err.Error()))
	}

	// Validate each repository
	for i, repo := range repositories {
		if err := validateRepository(repo); err != nil {
			return nil, errors.NewBadRequestError(fmt.Sprintf("invalid repository at index %d: %v", i, err.Error()))
		}
	}

	return repositories, nil
}

func validateRepository(repo *models.RepositoryModel) error {
	missingFields := []string{}

	if repo.Name == "" {
		missingFields = append(missingFields, "name")
	}
	if repo.Author == "" {
		missingFields = append(missingFields, "author")
	}
	if repo.ApiUrl == "" {
		missingFields = append(missingFields, "apiUrl")
	}
	if repo.CloneUrl == "" {
		missingFields = append(missingFields, "cloneUrl")
	}
	if repo.GhUrl == "" {
		missingFields = append(missingFields, "ghUrl")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
