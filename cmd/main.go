package main

import (
	"context"
	"log"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googlbye/internal/common"
	"github.com/jwtly10/googlbye/internal/models"
	"github.com/jwtly10/googlbye/internal/repository"
	"github.com/jwtly10/googlbye/internal/search"
	"go.uber.org/zap/zapcore"
)

func main() {
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config")
	}

	useJsonLogger := false
	logger := common.NewLogger(useJsonLogger, zapcore.DebugLevel)

	db, err := common.ConnectDB(config)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
	}

	repoRepo := repository.NewSQLRepoRepository(db)

	searchParams := &models.SearchParams{
		Query: "language:go stars:>1000",
		Opts: &github.SearchOptions{
			Sort:  "stars",
			Order: "desc",
			ListOptions: github.ListOptions{
				PerPage: 5,
			},
		},
		Page: 2,
	}

	search := search.NewRepoSearch(searchParams, config, logger, repoRepo)

	search.StartSearch(context.Background())
}
