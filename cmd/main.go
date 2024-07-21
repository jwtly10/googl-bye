package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/parser"
	"github.com/jwtly10/googl-bye/internal/repository"
	"github.com/jwtly10/googl-bye/internal/search"
	"go.uber.org/zap/zapcore"
)

func main() {
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config")
	}

	useJsonLogger := false
	logger := common.NewLogger(useJsonLogger, zapcore.InfoLevel)

	db, err := common.ConnectDB(config)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
	}

	repoRepo := repository.NewRepoRepository(db)
	stateRepo := repository.NewParserStateRepository(db)
	linkRepo := repository.NewParserLinkRepository(db)
	searchRepo := repository.NewSearchParamRepository(db)

	searchParams := &models.SearchParamsModel{
		Name:  "test_search",
		Query: "stars:300..1000 size:<=50000", // Repos with 300-1000 stars and smaller than 50MB
		Opts: github.SearchOptions{
			Sort:  "stars",
			Order: "desc",
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		},
		StartPage:      0,
		PagesToProcess: 5,
	}

	// Create a context that we can cancel to stop all goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start the search process asynchronously
	wg.Add(1)
	go func() {
		defer wg.Done()
		repoCache := search.NewRepoCache(repoRepo, logger)
		search := search.NewRepoSearch(searchParams, config, logger, repoRepo, repoCache, searchRepo)
		search.StartSearch(ctx)
	}()

	// The parser is a background job that should keep running in the background, to clear all items in the DB
	// It will try to run every 10 seconds
	// But will wait for the previous run to complete before continuing
	parser := parser.NewParser(logger, repoRepo, stateRepo, linkRepo)
	limit := 10

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		parser.StartParser(ctx, limit)
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				wg.Wait()

				wg.Add(1)
				go func() {
					defer wg.Done()
					parser.StartParser(ctx, limit)
				}()
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()

	wg.Wait()

	logger.Info("Parser has been gracefully shut down")
}
