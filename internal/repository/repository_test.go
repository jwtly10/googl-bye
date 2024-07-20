package repository_test

import (
	"context"
	"testing"

	"github.com/jwtly10/googlbye/internal/models"
	"github.com/jwtly10/googlbye/internal/repository"
	"github.com/jwtly10/googlbye/internal/test"
)

var (
	repos = []models.RepositoryModel{
		{
			Name:     "awesome-go",
			Author:   "avelino",
			ApiUrl:   "https://api.github.com/repos/avelino/awesome-go",
			GhUrl:    "https://github.com/avelino/awesome-go",
			CloneUrl: "https://github.com/avelino/awesome-go.git",
		},
		{
			Name:     "gin",
			Author:   "gin-gonic",
			ApiUrl:   "https://api.github.com/repos/gin-gonic/gin",
			GhUrl:    "https://github.com/gin-gonic/gin",
			CloneUrl: "https://github.com/gin-gonic/gin.git",
		},
	}
)

func TestRepoRepository_Integration(t *testing.T) {
	container, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())

	repoRepo := repository.NewSQLRepoRepository(db)

	t.Run("Create repos", func(t *testing.T) {
		for i := range repos {
			if err := repoRepo.CreateRepo(&repos[i]); err != nil {
				t.Errorf("expected no error when creating repo but got %v", err)
			}
			if repos[i].ID == "" {
				t.Error("expected repo ID to be set after creation")
			}
		}
	})

	t.Run("Update repo", func(t *testing.T) {
		repos[0].Name = "Updated Name"
		repos[0].Author = "Updated Author"
		if err := repoRepo.UpdateRepo(&repos[0]); err != nil {
			t.Errorf("expected no error when updating repo but got %v", err)
		}
		loaded, err := repoRepo.GetRepoByID(repos[0].ID)
		if err != nil {
			t.Errorf("expected no error when getting repo by id but got %v", err)
		}
		if loaded.Name != "Updated Name" {
			t.Errorf("expected loaded repo's name to be updated to 'Updated Name' but was '%s'", loaded.Name)
		}
		if loaded.Author != "Updated Author" {
			t.Errorf("expected loaded repo's author to be updated but was '%s'", loaded.Author)
		}
	})

	t.Run("Get repo by ID", func(t *testing.T) {
		loaded, err := repoRepo.GetRepoByID(repos[1].ID)
		if err != nil {
			t.Errorf("expected no error when getting repo by ID but got %v", err)
		}
		if loaded.ID != repos[1].ID {
			t.Errorf("expected loaded repo's id to match '%s' but was '%s'", repos[1].ID, loaded.ID)
		}
	})

	t.Run("Delete repos", func(t *testing.T) {
		for _, repo := range repos {
			if err := repoRepo.DeleteRepo(repo.ID); err != nil {
				t.Errorf("expected no error when deleting repo but got %v", err)
			}
			deletedRepo, err := repoRepo.GetRepoByID(repo.ID)
			if deletedRepo != nil {
				t.Errorf("expected deleted repo to be nil but was %v", deletedRepo)
			}
			if err != repository.ErrRepoNotFound {
				t.Errorf("expected ErrRepoNotFound when attempting to get repo by id after deletion, but got %v", err)
			}
		}
	})

	t.Run("Fail to update non existing repo", func(t *testing.T) {
		repos[0].Name = "Updated Name"
		repos[0].Author = "Updated Author"
		repos[0].ID = "bce72846-f5cc-4c3f-9d47-8af589dcd7bf"

		if err := repoRepo.UpdateRepo(&repos[0]); err == nil {
			t.Error("expected error while attempting to update a repo that should not exist")
		}
	})
}
