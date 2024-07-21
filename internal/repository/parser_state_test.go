package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
	"github.com/jwtly10/googl-bye/internal/test"
)

func TestParserStateRepository_Integration(t *testing.T) {
	container, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())

	parserStateRepo := repository.NewParserStateRepository(db)
	repoRepo := repository.NewRepoRepository(db)

	// Create a mock repo
	repo1 := models.RepositoryModel{
		Name:     "awesome-go",
		Author:   "avelino",
		ApiUrl:   "https://api.github.com/repos/avelino/awesome-go",
		GhUrl:    "https://github.com/avelino/awesome-go",
		CloneUrl: "https://github.com/avelino/awesome-go.git",
	}
	repo2 := models.RepositoryModel{
		Name:     "jwtly10",
		Author:   "googl-bye-test",
		ApiUrl:   "https://api.github.com/repos/jwtly10/googl-bye-test",
		GhUrl:    "https://github.com/avelino/googl-bye-test",
		CloneUrl: "https://github.com/avelino/googl-bye-test.git",
	}

	t.Run("Create repos", func(t *testing.T) {
		if err := repoRepo.CreateRepo(&repo1); err != nil {
			t.Errorf("expected no error when creating repo but got %v", err)
		}
		if repo1.ID == 0 {
			t.Error("expected repo ID to be set after creation")
		}
		if err := repoRepo.CreateRepo(&repo2); err != nil {
			t.Errorf("expected no error when creating repo but got %v", err)
		}
		if repo2.ID == 0 {
			t.Error("expected repo ID to be set after creation")
		}
	})

	t.Run("Set parser state", func(t *testing.T) {
		state := &models.ParserStateModel{
			Name:             "ParserJob",
			LastParsedRepoId: repo1.ID,
			LastParsedAt:     time.Now(),
		}

		err := parserStateRepo.SetParserState(state)
		if err != nil {
			t.Errorf("expected no error when setting parser state but got %v", err)
		}

		if state.ID == 0 {
			t.Error("expected parser state ID to be set after creation")
		}

		if state.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set")
		}

		if state.UpdatedAt.IsZero() {
			t.Error("expected UpdatedAt to be set")
		}
	})

	t.Run("Update existing parser state", func(t *testing.T) {
		initialState, err := parserStateRepo.GetParserState()
		if err != nil {
			t.Fatalf("expected no error when getting parser state but got %v", err)
		}

		// Hack to fix skew issues between golang time & postgres tz
		beforeUpdate := time.Now()

		updatedState := initialState
		updatedState.LastParsedRepoId = repo2.ID
		updatedState.LastParsedAt = time.Now().Add(time.Hour)

		err = parserStateRepo.SetParserState(updatedState)
		if err != nil {
			t.Fatalf("expected no error when updating parser state but got %v", err)
		}

		if updatedState.ID != initialState.ID {
			t.Errorf("expected ID to remain the same, got %d, want %d", updatedState.ID, initialState.ID)
		}

		if !updatedState.LastParsedAt.After(beforeUpdate) {
			t.Errorf("expected LastParsedAt to be later than before update. Got %v, want after %v",
				updatedState.LastParsedAt, beforeUpdate)
		}
	})
}
