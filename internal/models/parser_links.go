package models

import (
	"time"
)

// ParserLinksModel represents the repository data stored in the database.
type ParserLinksModel struct {
	Model
	RepoId      int    `db:"repo_id" json:"repoId"`
	Url         string `db:"url" json:"url"`
	ExpandedUrl string `db:"expanded_url" json:"expandedUrl"`
	File        string `db:"file" json:"file"`
	LineNumber  int    `db:"line_number" json:"lineNumber"`
	GithubUrl   string `db:"github_url" json:"github_url"`
	Path        string `db:"path" json:"path"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *ParserLinksModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
