package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/jwtly10/googl-bye/internal/models"
)

type RepoLinkRepository struct {
	db *sql.DB
}

func NewRepoLinkRepository(db *sql.DB) *RepoLinkRepository {
	return &RepoLinkRepository{db: db}
}

func (r *RepoLinkRepository) GetRepositoryWithLinks() ([]*models.RepoWithLinks, error) {
	rows, err := r.db.Query(`
        SELECT 
            r.id, r.name, r.author, r.state, r.api_url, r.gh_url, 
            r.language, r.stars, r.forks, r.size, r.last_push, r.clone_url, 
            r.error_msg, r.created_at, r.updated_at,
            l.id, l.url, l.expanded_url, l.file, l.line_number, l.github_url,
            l.path, l.created_at, l.updated_at
        FROM 
            public.repository_tb r
        LEFT JOIN 
            public.parser_links_tb l ON r.id = l.repo_id
        WHERE 
            (r.state = 'COMPLETED' OR r.state = 'ERROR')
            AND (r.state = 'ERROR' OR l.id IS NOT NULL)
        ORDER BY 
            r.id DESC, l.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	repositories := make(map[int]*models.RepoWithLinks)

	for rows.Next() {
		var r models.RepoWithLinks
		var l models.Link
		var linkID sql.NullInt64

		err := rows.Scan(
			&r.ID, &r.Name, &r.Author, &r.State, &r.ApiUrl, &r.GhUrl,
			&r.Language, &r.Stars, &r.Forks, &r.Size, &r.LastPush, &r.CloneURL,
			&r.ErrorMsg, &r.CreatedAt, &r.UpdatedAt,
			&linkID, &l.Url, &l.ExpandedURL, &l.File, &l.LineNumber, &l.GithubUrl,
			&l.Path, &l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		repo, exists := repositories[r.ID]
		if !exists {
			repo = &r
			repo.Links = make([]models.Link, 0)
			repositories[r.ID] = repo
		}

		if linkID.Valid {
			l.ID = int(linkID.Int64)
			repo.Links = append(repo.Links, l)
		}
	}

	result := make([]*models.RepoWithLinks, 0, len(repositories))
	for _, repo := range repositories {
		if repo.State == "COMPLETED" && len(repo.Links) == 0 {
			return nil, fmt.Errorf("repository with id %d has 'COMPLETED' status but no links", repo.ID)
		}
		result = append(result, repo)
	}

	return result, nil

}

func (r *RepoLinkRepository) GetRepositoryWithLinksForUser(author string) ([]*models.RepoWithLinks, error) {
	rows, err := r.db.Query(`
        SELECT 
            r.id, r.name, r.author, r.state, r.api_url, r.gh_url, 
            r.language, r.stars, r.forks, r.size, r.last_push, r.clone_url, 
            r.error_msg, r.created_at, r.updated_at,
            l.id, l.url, l.expanded_url, l.file, l.line_number, l.github_url,
            l.path, l.created_at, l.updated_at
        FROM 
            public.repository_tb r
        LEFT JOIN 
            public.parser_links_tb l ON r.id = l.repo_id
        WHERE 
            r.author = $1
        ORDER BY 
            r.id DESC, l.id
    `, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	repositories := make(map[int]*models.RepoWithLinks)
	for rows.Next() {
		var r models.RepoWithLinks
		var l models.Link
		var linkID, url, expandedURL, file, githubURL, path sql.NullString
		var linkCreatedAt, linkUpdatedAt sql.NullTime
		var lineNumber sql.NullInt16

		err := rows.Scan(
			&r.ID, &r.Name, &r.Author, &r.State, &r.ApiUrl, &r.GhUrl,
			&r.Language, &r.Stars, &r.Forks, &r.Size, &r.LastPush, &r.CloneURL,
			&r.ErrorMsg, &r.CreatedAt, &r.UpdatedAt,
			&linkID, &url, &expandedURL, &file, &lineNumber, &githubURL,
			&path, &linkCreatedAt, &linkUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		repo, exists := repositories[r.ID]
		if !exists {
			repo = &r
			repo.Links = make([]models.Link, 0)
			repositories[r.ID] = repo
		}

		if linkID.Valid {
			l.ID, _ = strconv.Atoi(linkID.String)
			l.Url = url.String
			l.ExpandedURL = expandedURL.String
			l.File = file.String
			l.LineNumber = int(lineNumber.Int16)
			l.GithubUrl = githubURL.String
			l.Path = path.String
			if linkCreatedAt.Valid {
				l.CreatedAt = linkCreatedAt.Time
			}
			if linkUpdatedAt.Valid {
				l.UpdatedAt = linkUpdatedAt.Time
			}
			repo.Links = append(repo.Links, l)
		}
	}

	result := make([]*models.RepoWithLinks, 0, len(repositories))
	for _, repo := range repositories {
		// if repo.State == "COMPLETED" && len(repo.Links) == 0 {
		// 	return nil, fmt.Errorf("repository with id %d has 'COMPLETED' status but no links", repo.ID)
		// }
		result = append(result, repo)
	}
	return result, nil
}
