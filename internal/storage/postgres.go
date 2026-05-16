package storage

import (
	"context"
	"time"

	"github.com/Dharshan2208/auth/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool becoz production servers don't open new db connection per request
// so we will reuse the pooled connections
type Store struct {
	DB *pgxpool.Pool
}

func New(databaseURL string) (*Store, error) {
	db, err := pgxpool.New(
		context.Background(),
		databaseURL,
	)
	if err != nil {
		return nil, err
	}

	return &Store{
		DB: db,
	}, nil
}

func (s *Store) CreateUser(user models.User) error {
	_, err := s.DB.Exec(
		context.Background(),
		`	
		INSERT INTO users (
			username,
			password_hash,
			role
		)
		VALUES ($1, $2, $3)
		`,
		user.Username,
		user.Password,
		user.Role,
	)

	return err
}

func (s *Store) GetUserByUsername(username string) (models.User, error) {
	var user models.User

	err := s.DB.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			username,
			password_hash,
			role,
			created_at
		FROM users
		WHERE username = $1
		`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	return user, err
}

func (s *Store) GetUserByID(id int) (models.User, error) {
	var user models.User

	err := s.DB.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			username,
			password_hash,
			role,
			created_at
		FROM users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	return user, err
}

func (s *Store) SaveRefreshToken(userID int, tokenHash string, expiresAt time.Time) error {
	_, err := s.DB.Exec(
		context.Background(),
		`
		INSERT INTO refresh_tokens
		(user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		`,
		userID,
		tokenHash,
		expiresAt,
	)

	return err
}

func (s *Store) GetUserIDByRefreshToken(tokenHash string) (int, error) {
	var userID int

	err := s.DB.QueryRow(
		context.Background(),
		`
		SELECT user_id
		FROM refresh_tokens
		WHERE token_hash = $1
		`,
		tokenHash,
	).Scan(&userID)

	return userID, err
}

func (s *Store) DeleteRefreshToken(
	tokenHash string,
) error {
	_, err := s.DB.Exec(
		context.Background(),
		`
		DELETE FROM refresh_tokens
		WHERE token_hash = $1
		`,
		tokenHash,
	)

	return err
}
