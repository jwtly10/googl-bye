package common

import (
	"fmt"
	"sync"

	"github.com/jwtly10/googl-bye/internal/repository"
)

type RepoCache struct {
	store sync.Map
}

func NewRepoCache(repoRepo repository.RepoRepository, log Logger) (*RepoCache, error) {
	log.Infof("Loading repo cache")

	cache := &RepoCache{}

	// Preload all repos
	repos, err := repoRepo.GetAllRepos()
	if err != nil {
		log.Errorf("Error getting all repos from db: %v", err)
		return nil, err
	}

	for _, repo := range repos {
		key := fmt.Sprintf("%s/%s", repo.Author, repo.Name)
		cache.store.Store(key, true)
	}

	log.Infof("Cache loaded with %d repos", len(repos))
	return cache, nil
}

// Set adds a key-value pair to the cache
func (c *RepoCache) Set(key, value interface{}) {
	c.store.Store(key, value)
}

// Get retrieves a value from the cache
func (c *RepoCache) Get(key interface{}) (interface{}, bool) {
	return c.store.Load(key)
}

// Exists checks if a key exists in the cache
func (c *RepoCache) Exists(key interface{}) bool {
	_, ok := c.store.Load(key)
	return ok
}

// Delete removes a key-value pair from the cache
func (c *RepoCache) Delete(key interface{}) {
	c.store.Delete(key)
}
