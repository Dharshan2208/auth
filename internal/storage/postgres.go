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

func (s *Store) CreateUser(ctx context.Context, user models.User) error {
	_, err := s.DB.Exec(
		ctx,
		`	
		INSERT INTO users (
			username,
			email,
			password_hash,
			role
		)
		VALUES ($1, $2, $3, $4)
		`,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
	)

	return err
}

func (s *Store) GetUserByUsernameOrEmail(ctx context.Context, identifier string) (models.User, error) {
	var user models.User

	err := s.DB.QueryRow(
		ctx,
		`
		SELECT
			id,
			username,
			email,
			password_hash,
			role,
			created_at
		FROM users
		WHERE username = $1 OR email = $1
		`,
		identifier,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	return user, err
}

func (s *Store) GetUserByID(ctx context.Context, id int) (models.User, error) {
	var user models.User

	err := s.DB.QueryRow(
		ctx,
		`
		SELECT
			id,
			username,
			email,
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
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	return user, err
}

func (s *Store) SaveRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := s.DB.Exec(
		ctx,
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

func (s *Store) GetUserIDByRefreshToken(ctx context.Context, tokenHash string) (int, error) {
	var userID int

	err := s.DB.QueryRow(
		ctx,
		`
		SELECT user_id
		FROM refresh_tokens
		WHERE token_hash = $1
		`,
		tokenHash,
	).Scan(&userID)

	return userID, err
}

// ConsumeRefreshToken atomically deletes the token and returns its user_id.
// If the token was already consumed, it returns an error...
// This eliminates the TOCTOU (Time-of-Check to Time-of-Use)race between checking and deleting.
func (s *Store) ConsumeRefreshToken(ctx context.Context, tokenHash string) (int, error) {
	var userID int

	err := s.DB.QueryRow(
		ctx,
		`
		DELETE FROM refresh_tokens
		WHERE token_hash = $1
		RETURNING user_id
		`,
		tokenHash,
	).Scan(&userID)

	return userID, err
}

func (s *Store) Ping(ctx context.Context) error {
	return s.DB.Ping(ctx)
}

func (s *Store) DeleteRefreshToken(
	ctx context.Context,
	tokenHash string,
) error {
	_, err := s.DB.Exec(
		ctx,
		`
		DELETE FROM refresh_tokens
		WHERE token_hash = $1
		`,
		tokenHash,
	)

	return err
}
