package repository

import (
	"context"
	"database/sql"
	"errors"
	"learn.oauth.client/data/model"
	"log"
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

func (repo *PostgresRepository) GetByAccessToken(accessToken string) (*model.TokenResponseData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
        SELECT access_token, token_type, expires_in, refresh_token, scope
        FROM oauth
        WHERE access_token = $1
    `

	var tokenData model.TokenResponseData
	err := repo.Conn.QueryRowContext(ctx, stmt, accessToken).Scan(
		&tokenData.AccessToken,
		&tokenData.TokenType,
		&tokenData.ExpiresIn,
		&tokenData.RefreshToken,
		&tokenData.Scope,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("TokenResponse not found by accessToken: %v\n", accessToken)
			return nil, nil
		}
		log.Printf("Error executing database query: %v\n", err)
		return nil, err
	}

	return &tokenData, nil
}
