package models

import (
	"time"
)

// ParserStateModel represents the repository data stored in the database.
type ParserStateModel struct {
	Model
	Name             string    `db:"name" json:"name"`
	LastParsedRepoId int       `db:"last_parsed_repo_id" json:"lastParsedRepoId"`
	LastParsedAt     time.Time `db:"last_parsed_at" json:"lastParsedAt"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *ParserStateModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
