package service

import (
	"fmt"
	"net/http"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/errors"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
)

type RepoLinkService struct {
	log common.Logger
	r   repository.RepoLinkRepository
}

func NewRepoLinkService(r repository.RepoLinkRepository, l common.Logger) *RepoLinkService {
	return &RepoLinkService{
		r:   r,
		log: l,
	}
}

func (rls *RepoLinkService) GetRepoLinks(r *http.Request) ([]*models.RepoWithLinks, error) {
	repoLinks, err := rls.r.GetRepositoryWithLinks()
	if err != nil {
		return nil, errors.NewInternalError(fmt.Sprintf("error when getting repo links: %v", err.Error()))
	}

	return repoLinks, nil
}

func (rls *RepoLinkService) GetUserRepoLinks(r *http.Request) ([]*models.RepoWithLinks, error) {
	username := r.URL.Query().Get("username")
	if username == "" {
		return nil, errors.NewBadRequestError("missing required field: username")
	}

	rls.log.Infof("Getting repo links for user: %s", username)

	repoLinks, err := rls.r.GetRepositoryWithLinksForUser(username)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Sprintf("error when getting repo links for user %v: %v", username, err.Error()))
	}

	return repoLinks, nil
}
