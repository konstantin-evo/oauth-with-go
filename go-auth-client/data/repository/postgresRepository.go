package repository

import (
	"context"
	"database/sql"
	"learn.oauth.client/data/model"
	"time"
)

const dbTimeout = time.Second * 3

type PostgresRepository struct {
	Conn *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		Conn: db,
	}
}

func (repo *PostgresRepository) Insert(token model.TokenResponseData) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		INSERT INTO oauth (access_token, token_type, expires_in, refresh_token, scope)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var newID int
	err := repo.Conn.QueryRowContext(ctx, stmt,
		token.AccessToken,
		token.TokenType,
		token.ExpiresIn,
		token.RefreshToken,
		token.Scope,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}
