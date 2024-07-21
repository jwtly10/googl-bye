package parser

import (
	"context"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/repository"
)

type Parser struct {
	rp       RepoParser
	log      common.Logger
	repoRepo repository.RepoRepository
}

func NewParser(rp RepoParser, log common.Logger, repoRepo repository.RepoRepository) *Parser {
	return &Parser{
		rp:       rp,
		log:      log,
		repoRepo: repoRepo,
	}
}

func (p *Parser) StartParser(ctx context.Context) {
	// 1. Find repos to parse

}
