package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/jwtly10/googlbye/internal/models"
)

type RepoRepository interface {
	CreateRepo(Repo *models.RepositoryModel) error
	GetRepoByID(id string) (*models.RepositoryModel, error)
	DeleteRepo(id string) error
	UpdateRepo(Repo *models.RepositoryModel) error
}

type sqlRepoRepository struct {
	database *sql.DB
}

func NewSQLRepoRepository(database *sql.DB) RepoRepository {
	return &sqlRepoRepository{database: database}
}

func (r *sqlRepoRepository) handleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRepoNotFound
	}
	if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
		return ErrRepoConnErr
	}
	return err
}

// CreateRepo inserts a new repo into the database
func (r *sqlRepoRepository) CreateRepo(repo *models.RepositoryModel) error {
	repo.BeforeCreate()
	query := `INSERT INTO public.repository_tb (name, author, api_url, gh_url, clone_url )
        VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.database.QueryRow(query,
		repo.Name,
		repo.Author,
		repo.ApiUrl,
		repo.GhUrl,
		repo.CloneUrl,
	).Scan(&repo.ID)
	if err != nil {
		return fmt.Errorf("failed to insert repo: %w", err)
	}
	repo.AfterCreate()
	return nil
}

// GetRepoByID retrieves a repo from the database by its unique ID
func (r *sqlRepoRepository) GetRepoByID(id string) (*models.RepositoryModel, error) {
	query := `SELECT id, name, author, api_url, gh_url, clone_url, created_at, updated_at FROM public.repository_tb WHERE id = $1`
	repo := &models.RepositoryModel{}
	err := r.database.QueryRow(query, id).Scan(
		&repo.ID,
		&repo.Name,
		&repo.Author,
		&repo.ApiUrl,
		&repo.GhUrl,
		&repo.CloneUrl,
		&repo.CreatedAt,
		&repo.UpdatedAt,
	)
	if err != nil {
		return nil, r.handleError(err)
	}
	return repo, nil
}

// UpdateRepo updates a repo in the database
func (r *sqlRepoRepository) UpdateRepo(repo *models.RepositoryModel) error {
	repo.BeforeUpdate()
	query := `UPDATE public.repository_tb SET name = $1, author = $2, api_url = $3, gh_url = $4, clone_url = $5 WHERE id = $6`
	if repo.CreatedAt.Unix() == 0 {
		return fmt.Errorf("unable to update a repo that was not loaded from the database")
	}
	rs, err := r.database.Exec(
		query,
		repo.Name,
		repo.Author,
		repo.ApiUrl,
		repo.GhUrl,
		repo.CloneUrl,
		repo.ID,
	)
	if err != nil {
		return r.handleError(err)
	}
	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrRepoNotFound
	}
	repo.AfterUpdate()
	return nil
}

// DeleteRepo deletes a repo from the database
func (r *sqlRepoRepository) DeleteRepo(id string) error {
	query := `DELETE FROM public.repository_tb WHERE id = $1`
	rs, err := r.database.Exec(query, id)
	if err != nil {
		return r.handleError(err)
	}
	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrRepoNotFound
	}
	return nil
}

var (
	ErrRepoNotFound = errors.New("repo not found") // ErrRepoNotFound is returned when a repo is not found in the database.
	ErrRepoConnErr  = errors.New("repository connection lost")
)
