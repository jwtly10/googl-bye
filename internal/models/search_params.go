package models

import "github.com/google/go-github/v39/github"

type SearchParams struct {
	Query string
	Opts  *github.SearchOptions
	Page  int
}
