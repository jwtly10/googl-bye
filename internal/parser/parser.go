package parser

import (
	"context"
	"fmt"
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

func (p *Parser) StartParser(ctx context.Context, limit int) {
	// 1. Find repos to parse
	reposToParse, err := p.repoRepo.GetPendingRepos()
	if err != nil {
		p.log.Errorf("Error getting pending repos: %v", err)
	}

	p.log.Infof("Found '%v' repos in state PENDING", len(reposToParse))

	// Parse the repos
	var lastRepo models.RepositoryModel
	for i, repo := range reposToParse {
		if i > limit {
			// We limit a 'run' to a certain number of repos
			p.log.Infof("Parse limit '%d' reached. Stopping parsing run.", limit)
			break
		}
		repo.ParseStatus = "PROCESSING"
		err = p.repoRepo.UpdateRepo(&repo)
		if err != nil {
			p.log.Errorf("[%s] Error updateing repo state : %v", fmt.Sprintf("%s/%s", repo.Author, repo.Name), err)
		}

		links, err := p.repoParser.ParseRepository(repo)
		if err != nil {
			p.log.Errorf("[%s] Error parsing repo: %v", fmt.Sprintf("%s/%s", repo.Author, repo.Name), err)
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
		repo.ParseStatus = "DONE"
		err = p.repoRepo.UpdateRepo(&repo)
		lastRepo = repo
	}

	if lastRepo != (models.RepositoryModel{}) {
		jobState, err := p.stateRepo.GetParserState()
		if err != nil {
			p.log.Errorf("[%s] Error getting job state: %v", fmt.Sprintf("%s/%s", lastRepo.Author, lastRepo.Name), err)
		}
		jobState.LastParsedAt = time.Now()
		jobState.LastParsedRepoId = lastRepo.ID
		err = p.stateRepo.SetParserState(jobState)
		if err != nil {
			p.log.Errorf("[%s] Error saving parser state: %v", fmt.Sprintf("%s/%s", lastRepo.Author, lastRepo.Name), err)
		}
	} else {
		p.log.Info("No last repo found during parsing. Will not update state.")
	}
}
