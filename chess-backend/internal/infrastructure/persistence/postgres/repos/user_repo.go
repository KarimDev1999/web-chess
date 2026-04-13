package repos

import (
	"context"
	"errors"

	"chess-backend/internal/domain/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	const query = `INSERT INTO users (id, email, password, username, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, u.ID, u.Email, u.Password, u.Username, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	const query = `SELECT id, email, password, username, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var u user.User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	const query = `SELECT id, email, password, username, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var u user.User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &u, err
}
