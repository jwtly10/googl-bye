package repository_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v39/github"
	"github.com/jwtly10/googl-bye/internal/models"
	"github.com/jwtly10/googl-bye/internal/repository"
	"github.com/jwtly10/googl-bye/internal/test"
)

func TestSearchParamRepository_Integration(t *testing.T) {
	container, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())

	searchParamRepo := repository.NewSearchParamRepository(db)

	t.Run("SaveSearchParams", func(t *testing.T) {
		// Define test data
		testOpts := github.SearchOptions{
			Sort:  "stars",
			Order: "desc",
		}

		params := &models.SearchParamsModel{
			Name:           "TestSearch",
			Query:          "golang",
			Opts:           testOpts,
			StartPage:      1,
			CurrentPage:    1,
			PagesToProcess: 5,
		}

		// Save the search params
		err := searchParamRepo.SaveSearchParams(params)
		if err != nil {
			t.Fatalf("Failed to save search params: %v", err)
		}

		// Check if ID was set
		if params.ID == 0 {
			t.Error("Expected ID to be set after saving, but it's still 0")
		}

		// Retrieve the saved params to verify
		retrievedParams, err := searchParamRepo.GetSearchParamsByID(params.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve saved search params: %v", err)
		}

		// Compare the retrieved params with the original
		if retrievedParams.Name != params.Name {
			t.Errorf("Expected Name to be %s, but got %s", params.Name, retrievedParams.Name)
		}
		if retrievedParams.Query != params.Query {
			t.Errorf("Expected Query to be %s, but got %s", params.Query, retrievedParams.Query)
		}
		if retrievedParams.StartPage != params.StartPage {
			t.Errorf("Expected StartPage to be %d, but got %d", params.StartPage, retrievedParams.StartPage)
		}
		if retrievedParams.CurrentPage != params.CurrentPage {
			t.Errorf("Expected CurrentPage to be %d, but got %d", params.CurrentPage, retrievedParams.CurrentPage)
		}
		if retrievedParams.PagesToProcess != params.PagesToProcess {
			t.Errorf("Expected PagesToProcess to be %d, but got %d", params.PagesToProcess, retrievedParams.PagesToProcess)
		}

		// Compare the Opts struct
		if retrievedParams.Opts.Sort != params.Opts.Sort {
			t.Errorf("Expected Opts.Sort to be %s, but got %s", params.Opts.Sort, retrievedParams.Opts.Sort)
		}
		if retrievedParams.Opts.Order != params.Opts.Order {
			t.Errorf("Expected Opts.Order to be %s, but got %s", params.Opts.Order, retrievedParams.Opts.Order)
		}

		// Test update
		params.Query = "updated golang"
		err = searchParamRepo.SaveSearchParams(params)
		if err != nil {
			t.Fatalf("Failed to update search params: %v", err)
		}

		// Retrieve updated params
		updatedParams, err := searchParamRepo.GetSearchParamsByID(params.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve updated search params: %v", err)
		}

		if updatedParams.Query != "updated golang" {
			t.Errorf("Expected updated Query to be 'updated golang', but got %s", updatedParams.Query)
		}
	})
}
