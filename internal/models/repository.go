package models

import (
	"time"
)

// RepositoryModel represents the repository data stored in the database.
type RepositoryModel struct {
	Model
	Name     string `db:"name" json:"name"`
	Author   string `db:"author" json:"author"`
	ApiUrl   string `db:"api_url" json:"apiUrl"`
	GhUrl    string `db:"gh_url" json:"ghUrl"`
	CloneUrl string `db:"clone_url" json:"cloneUrl"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *RepositoryModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
