package models

import "time"

type User struct {
	ID       int
	Username string `json:"username"`
	Email    string `json:"email"`
	// PasswordHash stores the bcrypt hash; never marshal it to JSON.
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	CreatedAt    time.Time
}
