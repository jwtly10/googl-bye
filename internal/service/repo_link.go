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
