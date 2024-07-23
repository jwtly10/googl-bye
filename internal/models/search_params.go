package models

import (
	"time"

	"github.com/google/go-github/v39/github"
)

// SearchParamsModel represents the repository data stored in the database.
type SearchParamsModel struct {
	Model
	Name           string               `db:"name" json:"name" validate:"omitempty,min=1"`
	Query          string               `db:"query" json:"query" validate:"omitempty,min=1"`
	Opts           github.SearchOptions `db:"opts" json:"opts"`
	StartPage      int                  `db:"start_page" json:"startPage" validate:"min=0"`
	CurrentPage    int                  `db:"current_page" json:"currentPage" validate:"min=0"`
	PagesToProcess int                  `db:"pages_to_process" json:"pagesToProcess" validate:"min=0"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *SearchParamsModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
