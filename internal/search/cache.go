package search

import (
	"fmt"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/repository"
)

func NewRepoCache(repoRepo repository.RepoRepository, log common.Logger) map[string]bool {
	repoCache := make(map[string]bool)

	// Preload all repos
	repos, err := repoRepo.GetAllRepos()
	if err != nil {
		log.Errorf("Error getting all repos from db")
	}

	for _, repo := range repos {
		repoCache[fmt.Sprintf("%s/%s", repo.Author, repo.Name)] = true
	}

	log.Info("Cache loaded")

	return repoCache
}
