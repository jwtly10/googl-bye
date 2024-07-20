package main

import (
	"context"
	"log"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googlbye/internal/common"
	"github.com/jwtly10/googlbye/internal/search"
	"go.uber.org/zap/zapcore"
)

func main() {

	config, err := common.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config")
	}

	useJsonLogger := false // TODO make configurable based on env
	logger := common.NewLogger(useJsonLogger, zapcore.DebugLevel)

	ghSearch := search.NewGithubSearch(config, logger)

	searchQuery := "language:go stars:>1000"
	opts := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: 5,
		},
	}

	ghSearch.FindRepositories(context.Background(), searchQuery, opts)
}
