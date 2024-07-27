package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"

	"github.com/jwtly10/googl-bye/internal/models"
)

type RepoLinkRepository interface {
	GetRepositoryWithLinks() ([]*models.RepoWithLinks, error)
	GetRepositoryWithLinksForUser(author string) ([]*models.RepoWithLinks, error)
	GetRepoLinksById(id int) (*models.RepoWithLinks, error)
}

type sqlRepoLinkRepository struct {
	database *sql.DB
}

func NewRepoLinkRepository(database *sql.DB) RepoLinkRepository {
	return &sqlRepoLinkRepository{database: database}
}

func (r *sqlRepoLinkRepository) handleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRepoNotFound
	}
	if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
		return ErrRepoConnErr
	}
	return err
}

func (r *sqlRepoLinkRepository) GetRepositoryWithLinks() ([]*models.RepoWithLinks, error) {
	rows, err := r.database.Query(`
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

func (r *sqlRepoLinkRepository) GetRepoLinksById(id int) (*models.RepoWithLinks, error) {
	query := `
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
            r.id = $1
        ORDER BY 
            r.id DESC, l.id
    `

	rows, err := r.database.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var repo *models.RepoWithLinks
	links := make([]models.Link, 0)

	for rows.Next() {
		if repo == nil {
			repo = &models.RepoWithLinks{}
		}

		var link models.Link
		var linkID, linkLineNumber sql.NullInt64
		var linkURL, linkExpandedURL, linkFile, linkGithubURL, linkPath sql.NullString
		var linkCreatedAt, linkUpdatedAt sql.NullTime

		err := rows.Scan(
			&repo.ID, &repo.Name, &repo.Author, &repo.State, &repo.ApiUrl, &repo.GhUrl,
			&repo.Language, &repo.Stars, &repo.Forks, &repo.Size, &repo.LastPush, &repo.CloneURL,
			&repo.ErrorMsg, &repo.CreatedAt, &repo.UpdatedAt,
			&linkID, &linkURL, &linkExpandedURL, &linkFile, &linkLineNumber, &linkGithubURL,
			&linkPath, &linkCreatedAt, &linkUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if linkID.Valid {
			link.ID = int(linkID.Int64)
			link.Url = linkURL.String
			link.ExpandedURL = linkExpandedURL.String
			link.File = linkFile.String
			link.LineNumber = int(linkLineNumber.Int64)
			link.GithubUrl = linkGithubURL.String
			link.Path = linkPath.String
			link.CreatedAt = linkCreatedAt.Time
			link.UpdatedAt = linkUpdatedAt.Time
			links = append(links, link)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if repo == nil {
		return nil, fmt.Errorf("repository with ID %d not found", id)
	}

	repo.Links = links
	return repo, nil
}

func (r *sqlRepoLinkRepository) GetRepositoryWithLinksForUser(author string) ([]*models.RepoWithLinks, error) {
	rows, err := r.database.Query(`
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
