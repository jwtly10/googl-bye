package models

import "time"

// RepoWithLinks  is a DTO for the frontend.
// Combines models.RepositoryModel and models.ParseLinksModel, to prevent additional frontend parsing logic that will need to be maintained
type RepoWithLinks struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	State     string    `json:"state"`
	ApiUrl    string    `json:"apiUrl"`
	GhUrl     string    `json:"ghUrl"`
	Language  string    `json:"language"`
	Stars     int       `json:"stars"`
	Forks     int       `json:"forks"`
	Size      int       `json:"size"`
	LastPush  time.Time `json:"lastPush"`
	CloneURL  string    `json:"cloneUrl"`
	ErrorMsg  string    `json:"errorMsg"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Links     []Link    `json:"links"`
}

type Link struct {
	ID          int       `json:"id"`
	Url         string    `json:"url"`
	ExpandedURL string    `json:"expandedUrl"`
	File        string    `json:"file"`
	LineNumber  int       `json:"lineNumber"`
	GithubUrl   string    `json:"githubUrl"`
	Path        string    `json:"path"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
