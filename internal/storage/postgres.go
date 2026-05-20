package storage

import (
	"context"
	"time"

	"github.com/Dharshan2208/auth/internal/models"
	"github.com/jackc/pgx/v5"
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

func (s *Store) CreateSession(ctx context.Context, userID int, refreshHash string, expiresAt time.Time, ip string, userAgent string) error {
	_, err := s.DB.Exec(
		ctx,
		`
		INSERT INTO sessions
			(user_id, refresh_hash, expires_at, ip, user_agent)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''))
		`,
		userID,
		refreshHash,
		expiresAt,
		ip,
		userAgent,
	)

	return err
}

// RotateSession atomically replaces the refresh hash and returns the owning user_id.
// Only non-revoked, non-expired sessions can be rotated.
func (s *Store) RotateSession(ctx context.Context, oldRefreshHash string, newRefreshHash string, newExpiresAt time.Time, now time.Time, ip string, userAgent string) (int, error) {
	var userID int

	err := s.DB.QueryRow(
		ctx,
		`
		UPDATE sessions
		SET
			refresh_hash = $2,
			expires_at = $3,
			last_used_at = $4,
			ip = NULLIF($5, ''),
			user_agent = NULLIF($6, '')
		WHERE
			refresh_hash = $1
			AND revoked_at IS NULL
			AND expires_at > $4
		RETURNING user_id
		`,
		oldRefreshHash,
		newRefreshHash,
		newExpiresAt,
		now,
		ip,
		userAgent,
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
	cmd, err := s.DB.Exec(
		ctx,
		`
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE refresh_hash = $1 AND revoked_at IS NULL
		`,
		tokenHash,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// RevokeAllSessionsForUserDevice revokes all sessions for a user on a given device.
func (s *Store) RevokeAllSessionsForUserDevice(ctx context.Context, userID int, userAgent string) error {
	_, err := s.DB.Exec(
		ctx,
		`
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE user_id = $1
		  AND revoked_at IS NULL
		  AND (
			($2 <> '' AND user_agent = $2)
			OR ($2 = '')
		  )
		`,
		userID,
		userAgent,
	)
	return err
}

func (s *Store) RevokeAllSessionsForUser(ctx context.Context, userID int) error {
	_, err := s.DB.Exec(
		ctx,
		`
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
		`,
		userID,
	)
	return err
}

func (s *Store) UpdateUserPasswordHash(ctx context.Context, userID int, passwordHash string) error {
	cmd, err := s.DB.Exec(
		ctx,
		`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1
		`,
		userID,
		passwordHash,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
