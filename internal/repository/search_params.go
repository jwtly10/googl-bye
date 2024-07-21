package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jwtly10/googl-bye/internal/models"
)

type SearchParamRepository interface {
	SaveSearchParams(params *models.SearchParamsModel) error
	GetSearchParamsByID(id int) (*models.SearchParamsModel, error)
	GetSearchParamsByName(name string) (*models.SearchParamsModel, error)
}
type sqlSearchParamRepository struct {
	database *sql.DB
}

func NewSearchParamRepository(database *sql.DB) SearchParamRepository {
	return &sqlSearchParamRepository{database: database}
}

func (r *sqlSearchParamRepository) SaveSearchParams(params *models.SearchParamsModel) error {
	params.BeforeCreate()

	// Storing the opts as json
	jsonOpts, err := json.Marshal(params.Opts)
	if err != nil {
		return fmt.Errorf("failed marshal opts: %w", err)
	}

	query := `
		INSERT INTO public.search_params_history_tb (name, query, opts, start_page, current_page, pages_to_process)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (name) DO UPDATE
		SET query = EXCLUDED.query,
			opts = EXCLUDED.opts,
			start_page = EXCLUDED.start_page,
			current_page = EXCLUDED.current_page,
			pages_to_process = EXCLUDED.pages_to_process,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, name, query, opts, start_page, current_page, pages_to_process`

	var optsJSON []byte
	err = r.database.QueryRow(query,
		params.Name,
		params.Query,
		jsonOpts,
		params.StartPage,
		params.CurrentPage,
		params.PagesToProcess,
	).Scan(&params.ID,
		&params.Name,
		&params.Query,
		&optsJSON,
		&params.StartPage,
		&params.CurrentPage,
		&params.PagesToProcess,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert parser state: %w", err)
	}

	// Unmarshal the json stored into the struct
	err = json.Unmarshal(optsJSON, &params.Opts)
	if err != nil {
		return fmt.Errorf("failed to unmarshal opts: %w", err)
	}

	params.AfterCreate()
	return nil
}

func (r *sqlSearchParamRepository) GetSearchParamsByID(id int) (*models.SearchParamsModel, error) {
	query := `
		SELECT id, name, query, opts, start_page, current_page, pages_to_process, created_at, updated_at
		FROM public.search_params_history_tb
		WHERE id = $1`

	var params models.SearchParamsModel
	var optsJSON []byte

	err := r.database.QueryRow(query, id).Scan(
		&params.ID,
		&params.Name,
		&params.Query,
		&optsJSON,
		&params.StartPage,
		&params.CurrentPage,
		&params.PagesToProcess,
		&params.CreatedAt,
		&params.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("search params with id %d not found", id)
		}
		return nil, fmt.Errorf("error querying search params: %w", err)
	}

	// Unmarshal the JSON-encoded opts back into the struct
	err = json.Unmarshal(optsJSON, &params.Opts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal opts: %w", err)
	}

	return &params, nil
}

func (r *sqlSearchParamRepository) GetSearchParamsByName(name string) (*models.SearchParamsModel, error) {
	query := `
		SELECT id, name, query, opts, start_page, current_page, pages_to_process, created_at, updated_at
		FROM public.search_params_history_tb
		WHERE name = $1`

	var params models.SearchParamsModel
	var optsJSON []byte

	err := r.database.QueryRow(query, name).Scan(
		&params.ID,
		&params.Name,
		&params.Query,
		&optsJSON,
		&params.StartPage,
		&params.CurrentPage,
		&params.PagesToProcess,
		&params.CreatedAt,
		&params.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying search params: %w", err)
	}

	// Unmarshal the JSON-encoded opts back into the struct
	err = json.Unmarshal(optsJSON, &params.Opts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal opts: %w", err)
	}

	return &params, nil
}
