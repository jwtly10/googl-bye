package parser

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
)

type Parser struct {
	repoParser RepoParser
	log        common.Logger
	repoRepo   repository.RepoRepository
	stateRepo  repository.ParserStateRepository
	linkRepo   repository.ParserLinksRepository
}

func NewParser(log common.Logger, repoRepo repository.RepoRepository, stateRepo repository.ParserStateRepository, linkRepo repository.ParserLinksRepository) *Parser {
	git := NewGitCmdLine(log)
	rp := NewRepoParser(git, log)

	return &Parser{
		repoParser: *rp,
		log:        log,
		repoRepo:   repoRepo,
		linkRepo:   linkRepo,
		stateRepo:  stateRepo,
	}
}

// StartParser finds repositories that are due to be parsed (status PENDING)
// It will pull 'limit' repos from DB and process them asynchronously
// Parsing cannot take more than 30 seconds
func (p *Parser) StartParser(ctx context.Context, limit int) {
	p.log.Info("Starting parser run")

	if limit == -1 {
		p.log.Info("Limit set to -1. Not parsing.")
		return
	}

	reposToParse, err := p.repoRepo.GetPendingRepos()
	if err != nil {
		p.log.Errorf("Error getting pending repos: %v", err)
	}

	p.log.Infof("Found '%v' repos in state PENDING", len(reposToParse))

	var wg sync.WaitGroup
	resultChan := make(chan models.RepositoryModel, limit)

	for i, repo := range reposToParse {
		if i > limit {
			// We limit a 'run' to a certain number of repos
			p.log.Infof("Parse limit '%d' reached. Stopping parsing run.", limit)
			break
		}

		wg.Add(1)

		go func(repo models.RepositoryModel) {
			defer wg.Done()

			timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			done := make(chan struct{})

			go func() {
				defer close(done)
				repo.State = "PROCESSING"
				err = p.repoRepo.UpdateRepo(&repo)
				if err != nil {
					p.log.Errorf("[%s] Error updateing repo state : %v", fmt.Sprintf("%s/%s", repo.Author, repo.Name), err)
				}

				links, err := p.repoParser.ParseRepository(repo)
				if err != nil {
					p.log.Errorf("[%s] Error parsing repo: %v", fmt.Sprintf("%s/%s", repo.Author, repo.Name), err)
					// If this fails, we should set state failed
					repo.State = "ERROR"
					repo.ErrorMsg = err.Error()
					err = p.repoRepo.UpdateRepo(&repo)
					if err != nil {
						p.log.Errorf("[%s] Error updating repo state: %v", fmt.Sprintf("%s/%s", repo.Author, repo.Name), err)
					}
				}

				// Save any links
				p.log.Infof("[%s] Found '%v' goo.gl links", fmt.Sprintf("%s/%s", repo.Author, repo.Name), len(links))
				for _, link := range links {
					link.RepoId = repo.ID
					err = p.linkRepo.CreateParserLink(&link)
					if err != nil {
						p.log.Errorf("[%s] Error saving repo link: %v", fmt.Sprintf("%s/%s", repo.Author, repo.Name), err)

					}
				}

				// Update states on success
				repo.State = "COMPLETED"
				err = p.repoRepo.UpdateRepo(&repo)
				if err != nil {
					p.log.Errorf("[%s] Error updating repo state: %v", fmt.Sprintf("%s/%s", repo.Author, repo.Name), err)
				}

				resultChan <- repo

			}()

			select {
			case <-timeoutCtx.Done():
				if timeoutCtx.Err() == context.DeadlineExceeded {
					p.log.Warnf("[%s] Processing timed out after 30 seconds", fmt.Sprintf("%s/%s", repo.Author, repo.Name))
					repo.State = "TIMEOUT"
					err := p.repoRepo.UpdateRepo(&repo)
					if err != nil {
						p.log.Errorf("[%s] Error updating repo state after timeout: %v", fmt.Sprintf("%s/%s", repo.Author, repo.Name), err)
					}
				}
			case <-done:
				// Processing completed within the timeout
			}

		}(repo)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var lastRepo models.RepositoryModel
	for repo := range resultChan {
		lastRepo = repo
	}

	if lastRepo != (models.RepositoryModel{}) {
		jobState, err := p.stateRepo.GetParserState()
		if err != nil {
			p.log.Errorf("[%s] Error getting job state: %v", fmt.Sprintf("%s/%s", lastRepo.Author, lastRepo.Name), err)
		}
		jobState.Name = "parser_job"
		jobState.LastParsedAt = time.Now()
		jobState.LastParsedRepoId = lastRepo.ID
		err = p.stateRepo.SetParserState(jobState)
		if err != nil {
			p.log.Errorf("[%s] Error saving parser state: %v", fmt.Sprintf("%s/%s", lastRepo.Author, lastRepo.Name), err)
		}
	}
}
