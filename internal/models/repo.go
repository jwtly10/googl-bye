package models

import (
	"time"
)

// RepoModel represents the repository data stored in the database.
type RepoModel struct {
	Model
	Name     string `db:"name" json:"name"`
	Author   string `db:"author" json:"author"`
	ApiUrl   string `db:"url" json:"url"`
	GhUrl    string `db:"ghUrl" json:"ghUrl"`
	CloneUrl string `db:"cloneUrl" json:"cloneUrl"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *RepoModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
