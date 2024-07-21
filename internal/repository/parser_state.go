package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/jwtly10/googl-bye/internal/models"
)

type ParserStateRepository interface {
	SetParserState(Repo *models.ParserStateModel) error
	GetParserState() (*models.ParserStateModel, error)
}

type sqlParserStateRepository struct {
	database *sql.DB
}

func NewParserStateRepository(database *sql.DB) ParserStateRepository {
	return &sqlParserStateRepository{database: database}
}

func (r *sqlParserStateRepository) handleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRepoNotFound
	}
	if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
		return ErrRepoConnErr
	}
	return err
}

// TODO FIX THIS!!!!
// GetParserState gets the current state of the parser
func (r *sqlParserStateRepository) GetParserState() (*models.ParserStateModel, error) {
	id := 1
	query := `SELECT id, last_parsed_repo_id, last_parsed_at, created_at, updated_at FROM public.parser_state_tb WHERE id = $1`
	state := &models.ParserStateModel{}
	err := r.database.QueryRow(query, id).Scan(
		&state.ID,
		&state.LastParsedRepoId,
		&state.LastParsedAt,
		&state.CreatedAt,
		&state.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Nicely handly if there is no state yet
			return &models.ParserStateModel{}, nil
		}
		return nil, r.handleError(err)
	}
	return state, nil
}

// SetParserState inserts or updates the current parser state
func (r *sqlParserStateRepository) SetParserState(state *models.ParserStateModel) error {
	state.BeforeCreate()
	query := `
		INSERT INTO public.parser_state_tb (last_parsed_repo_id, last_parsed_at)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET last_parsed_repo_id = EXCLUDED.last_parsed_repo_id,
			last_parsed_at = EXCLUDED.last_parsed_at,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, last_parsed_at, created_at, updated_at`

	err := r.database.QueryRow(query,
		state.LastParsedRepoId,
		state.LastParsedAt,
	).Scan(&state.ID, &state.LastParsedAt, &state.CreatedAt, &state.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to upsert parser state: %w", err)
	}

	state.AfterCreate()
	return nil
}
