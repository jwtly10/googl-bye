package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jwtly10/googl-bye/api"
	"github.com/jwtly10/googl-bye/api/handlers"
	"github.com/jwtly10/googl-bye/api/middleware"
	"github.com/jwtly10/googl-bye/api/routes"
	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/parser"
	"github.com/jwtly10/googl-bye/internal/repository"
	"github.com/jwtly10/googl-bye/internal/search"
	"github.com/jwtly10/googl-bye/internal/service"
	"go.uber.org/zap/zapcore"
)

func main() {

	// ***** APPLICATION SETUP *****

	// Setup config
	config, err := common.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	useJsonLogger := false
	logger := common.NewLogger(useJsonLogger, zapcore.DebugLevel)

	// Setup db
	db, err := common.ConnectDB(config)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
	}

	// Init repos
	repoRepo := repository.NewRepoRepository(db)
	repoLinkRepo := repository.NewRepoLinkRepository(db)
	stateRepo := repository.NewParserStateRepository(db)
	linkRepo := repository.NewParserLinkRepository(db)
	searchRepo := repository.NewSearchParamRepository(db)

	// Init repo cache
	repoCache, err := common.NewRepoCache(repoRepo, logger)
	if err != nil {
		logger.Fatalf("Failed to create a repo cache on init: %v", err)
	}

	// ***** SERVER SETUP *****

	// Setup main router
	router := api.NewAppRouter(logger)

	// Attach global middleware
	loggerMw := middleware.NewRequestLoggerMiddleware(logger)

	// Setup frontend routes
	// router.SetupSwagger() // TODO
	router.ServeStaticFiles("./react/dist") // TODO: Should this only happen in dev. In prod we package binary with the frontend, and just ship the binary

	// Setup Github route
	ghs := search.NewGithubSearch(config, logger, searchRepo, repoRepo)
	githubService := service.NewGithubService(*ghs, logger)
	githubHandler := handlers.NewGithubHandler(logger, *githubService)
	routes.NewGithubRoutes(router, logger, *githubHandler, loggerMw)

	// Setup Repo route
	repoService := service.NewRepoService(repoRepo, logger, repoCache)
	repoHandler := handlers.NewRepoHandler(logger, *repoService)
	routes.NewRepoRoutes(router, logger, *repoHandler, loggerMw)

	// Setup RepoLink route
	repoLinkService := service.NewRepoLinkService(*repoLinkRepo, logger)
	repoLinkHandler := handlers.NewRepoLinkHandler(logger, *repoLinkService)
	routes.NewRepoLinkRoutes(router, logger, *repoLinkHandler, loggerMw)

	// Create a context that we can cancel to stop all goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting server", err)
		}
	}()

	// Start parser
	parser := parser.NewParser(logger, repoRepo, stateRepo, linkRepo)
	limit := 10
	ticker := time.NewTicker(time.Duration(config.ParserInterval) * time.Second)
	logger.Infof("Parser Job running every '%d' seconds", config.ParserInterval)
	defer ticker.Stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				parser.StartParser(ctx, limit)
			}
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Initiate graceful shutdown
	logger.Info("Initiating graceful shutdown...")

	// Cancel the context to stop background jobs
	cancel()

	// Shutdown the HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error shutting down server", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	logger.Info("Server and background jobs have been gracefully shut down")
}
