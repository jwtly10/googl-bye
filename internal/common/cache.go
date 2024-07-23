package common

import (
	"fmt"
	"sync"

	"github.com/jwtly10/googl-bye/internal/repository"
)

func NewRepoCache(repoRepo repository.RepoRepository, log Logger) *sync.Map {
	repoCache := &sync.Map{}

	// Preload all repos
	repos, err := repoRepo.GetAllRepos()
	if err != nil {
		log.Errorf("Error getting all repos from db: %v", err)
		return repoCache
	}

	for _, repo := range repos {
		key := fmt.Sprintf("%s/%s", repo.Author, repo.Name)
		repoCache.Store(key, true)
	}

	log.Infof("Cache loaded with %d repos", len(repos))
	return repoCache
}
