package repository_test

import (
	"context"
	"testing"

	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
	"github.com/jwtly10/googl-bye/internal/test"
)

var (
	parserLinks = []models.ParserLinksModel{
		{
			RepoId:      1,
			Url:         "https://example.com",
			ExpandedUrl: "https://www.example.com",
			File:        "README.md",
			LineNumber:  10,
			Path:        "/docs/README.md",
		},
		{
			RepoId:      1,
			Url:         "https://google.com",
			ExpandedUrl: "https://www.google.com",
			File:        "main.go",
			LineNumber:  25,
			Path:        "/src/main.go",
		},
	}
)

func TestParserLinkRepository_Integration(t *testing.T) {
	container, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())

	parserLinkRepo := repository.NewParserLinkRepository(db)
	repoRepo := repository.NewRepoRepository(db)

	// Create a mock repo
	repo := models.RepositoryModel{
		Name:     "awesome-go",
		Author:   "avelino",
		ApiUrl:   "https://api.github.com/repos/avelino/awesome-go",
		GhUrl:    "https://github.com/avelino/awesome-go",
		CloneUrl: "https://github.com/avelino/awesome-go.git",
	}

	t.Run("Create repo", func(t *testing.T) {
		if err := repoRepo.CreateRepo(&repo); err != nil {
			t.Errorf("expected no error when creating repo but got %v", err)
		}
		if repo.ID == 0 {
			t.Error("expected repo ID to be set after creation")
		}
	})

	t.Run("Create parser links", func(t *testing.T) {
		for i := range parserLinks {
			parserLinks[i].RepoId = repo.ID
			if err := parserLinkRepo.CreateParserLink(&parserLinks[i]); err != nil {
				t.Errorf("expected no error when creating parser link but got %v", err)
			}
			if parserLinks[i].ID == 0 {
				t.Error("expected parser link ID to be set after creation")
			}
		}
	})

	t.Run("Error when duplicate parser link", func(t *testing.T) {
		duplicateLink := parserLinks[0]
		duplicateLink.ID = 0 // Reset ID to simulate a new insertion
		err := parserLinkRepo.CreateParserLink(&duplicateLink)
		if err == nil {
			t.Error("expected error when creating duplicate parser link, but got nil")
		}
	})
}
