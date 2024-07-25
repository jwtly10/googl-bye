package models

import (
	"time"
)

// RepositoryModel represents the repository data stored in the database.
type RepositoryModel struct {
	Model
	Name     string    `db:"name" json:"name"`
	Author   string    `db:"author" json:"author"`
	State    string    `db:"state" json:"state"`
	Language string    `db:"language" json:"language"`
	Stars    int       `db:"stars" json:"stars"`
	Forks    int       `db:"forks" json:"forks"`
	Size     int       `db:"size" json:"size"`
	LastPush time.Time `db:"last_push" json:"lastPush"`
	ApiUrl   string    `db:"api_url" json:"apiUrl"`
	GhUrl    string    `db:"gh_url" json:"ghUrl"`
	CloneUrl string    `db:"clone_url" json:"cloneUrl"`
	ErrorMsg string    `db:"error_msg" json:"essorMsg"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *RepositoryModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
