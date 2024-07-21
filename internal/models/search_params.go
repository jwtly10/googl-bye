package models

import (
	"time"

	"github.com/google/go-github/v39/github"
)

// SearchParamsModel represents the repository data stored in the database.
type SearchParamsModel struct {
	Model
	Name           string               `db:"name" json:"name"`
	Query          string               `db:"query" json:"query"`
	Opts           github.SearchOptions `db:"opts" json:"opts"`
	StartPage      int                  `db:"start_page" json:"startPage"`
	CurrentPage    int                  `db:"current_page" json:"currentPage"`
	PagesToProcess int                  `db:"pages_to_process" json:"pagesToProcess"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *SearchParamsModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
