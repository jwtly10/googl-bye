package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/jwtly10/googl-bye/internal/models"
)

type ParserLinksRepository interface {
	CreateParserLink(Repo *models.ParserLinksModel) error
}

type sqlParserLinkRepository struct {
	database *sql.DB
}

func NewParserLinkRepository(database *sql.DB) ParserLinksRepository {
	return &sqlParserLinkRepository{database: database}
}

func (r *sqlParserLinkRepository) handleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRepoNotFound
	}
	if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
		return ErrRepoConnErr
	}
	return err
}

// CreateParserLink inserts a new link into the database
func (r *sqlParserLinkRepository) CreateParserLink(link *models.ParserLinksModel) error {
	link.BeforeCreate()
	query := `INSERT INTO public.parser_links_tb (repo_id, url, expanded_url, file, line_number, github_url, path)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.database.QueryRow(query,
		link.RepoId,
		link.Url,
		link.ExpandedUrl,
		link.File,
		link.LineNumber,
		link.GithubUrl,
		link.Path,
	).Scan(&link.ID)
	if err != nil {
		return fmt.Errorf("failed to insert link: %w", err)
	}
	link.AfterCreate()
	return nil
}
