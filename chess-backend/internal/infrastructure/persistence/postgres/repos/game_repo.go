package repos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"chess-backend/internal/domain/chess"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameRepository struct {
	db *pgxpool.Pool
}

func NewGameRepository(db *pgxpool.Pool) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) Save(ctx context.Context, game *chess.Game) error {
	movesJSON, err := json.Marshal(game.Moves)
	if err != nil {
		return err
	}

	var result interface{}
	if game.Result != nil {
		result = string(*game.Result)
	}

	var endReason interface{}
	if game.EndReason != nil {
		endReason = string(*game.EndReason)
	}

	const query = `INSERT INTO games (
		id, white_player_id, black_player_id, status, fen, turn, moves,
		result, end_reason, created_at, updated_at,
		time_base, time_increment, white_remaining, black_remaining, last_move_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err = r.db.Exec(ctx, query,
		game.ID, game.WhitePlayerID, game.BlackPlayerID,
		game.Status, game.CurrentFEN, game.Turn, movesJSON,
		result, endReason, game.CreatedAt, game.UpdatedAt,
		game.TimeControl.Base, game.TimeControl.Increment,
		game.WhiteRemaining, game.BlackRemaining, game.LastMoveAt,
	)
	return err
}

func (r *GameRepository) FindByID(ctx context.Context, id chess.GameID) (*chess.Game, error) {
	const query = `SELECT
		id, white_player_id, black_player_id, status, fen, turn, moves,
		result, end_reason, created_at, updated_at,
		time_base, time_increment, white_remaining, black_remaining, last_move_at
		FROM games WHERE id = $1`
	game, err := r.scanGameRow(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return game, err
}

func (r *GameRepository) FindWaitingGames(ctx context.Context) ([]*chess.Game, error) {
	const query = `SELECT
		id, white_player_id, black_player_id, status, fen, turn, moves,
		result, end_reason, created_at, updated_at,
		time_base, time_increment, white_remaining, black_remaining, last_move_at
		FROM games WHERE status = 'waiting'`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*chess.Game
	for rows.Next() {
		game, err := r.scanGameRow(rows)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, rows.Err()
}

func (r *GameRepository) FindByPlayerID(ctx context.Context, playerID string) ([]*chess.Game, error) {
	const query = `SELECT
		id, white_player_id, black_player_id, status, fen, turn, moves,
		result, end_reason, created_at, updated_at,
		time_base, time_increment, white_remaining, black_remaining, last_move_at
		FROM games WHERE white_player_id = $1 OR black_player_id = $1 ORDER BY updated_at DESC`
	rows, err := r.db.Query(ctx, query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*chess.Game
	for rows.Next() {
		game, err := r.scanGameRow(rows)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, rows.Err()
}

func (r *GameRepository) Update(ctx context.Context, game *chess.Game) error {
	movesJSON, err := json.Marshal(game.Moves)
	if err != nil {
		return err
	}

	var result interface{}
	if game.Result != nil {
		result = string(*game.Result)
	} else {
		result = nil
	}

	var endReason interface{}
	if game.EndReason != nil {
		endReason = string(*game.EndReason)
	} else {
		endReason = nil
	}

	const query = `UPDATE games SET
		white_player_id=$2, black_player_id=$3, status=$4, fen=$5, turn=$6,
		moves=$7, result=$8, end_reason=$9, updated_at=$10,
		time_base=$11, time_increment=$12,
		white_remaining=$13, black_remaining=$14, last_move_at=$15
		WHERE id=$1`
	_, err = r.db.Exec(ctx, query,
		game.ID, game.WhitePlayerID, game.BlackPlayerID,
		game.Status, game.CurrentFEN, game.Turn, movesJSON,
		result, endReason, game.UpdatedAt,
		game.TimeControl.Base, game.TimeControl.Increment,
		game.WhiteRemaining, game.BlackRemaining, game.LastMoveAt,
	)
	return err
}

func (r *GameRepository) scanGameRow(row interface{ Scan(...interface{}) error }) (*chess.Game, error) {
	var game chess.Game
	var movesJSON []byte

	err := row.Scan(
		&game.ID, &game.WhitePlayerID, &game.BlackPlayerID, &game.Status,
		&game.CurrentFEN, &game.Turn, &movesJSON, &game.Result, &game.EndReason,
		&game.CreatedAt, &game.UpdatedAt,
		&game.TimeControl.Base, &game.TimeControl.Increment,
		&game.WhiteRemaining, &game.BlackRemaining, &game.LastMoveAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan game row: %w", err)
	}

	if err := json.Unmarshal(movesJSON, &game.Moves); err != nil {
		return nil, err
	}
	return &game, nil
}
