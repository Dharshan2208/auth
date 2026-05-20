package handlers

// Request Types

// SignupRequest represents the request body for user registration.
type SignupRequest struct {
	Username string `json:"username" example:"johndoe" validate:"required,min=3,max=31"`
	Email    string `json:"email" example:"john@example.com" validate:"required,email"`
	Password string `json:"password" example:"Str0ng!Pass1" validate:"required,min=12,max=72"`
}

// LoginRequest represents the request body for login.
// The identifier can be either a username or an email.
type LoginRequest struct {
	Identifier string `json:"identifier" example:"johndoe"`
	Password   string `json:"password" example:"Str0ng!Pass1"`
}

// RefreshTokenRequest represents the request body for token refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// LogoutRequest represents the request body for logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// ChangePasswordRequest represents the request body for changing a password.
type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" example:"Str0ng!Pass1"`
	ConfirmOldPassword string `json:"confirm_old_password" example:"Str0ng!Pass1"`
	NewPassword        string `json:"new_password" example:"NewStr0ng!Pass2"`
}

// Response Types

// MessageResponse is a generic response with a message string.
type MessageResponse struct {
	Message string `json:"message" example:"operation successful"`
}

// ErrorResponse is the standard error response body.
type ErrorResponse struct {
	Error string `json:"error" example:"description of what went wrong"`
}

// TokenPairResponse contains both access and refresh tokens.
type TokenPairResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// ProfileResponse contains the authenticated user's profile data.
type ProfileResponse struct {
	UserID   int    `json:"user_id" example:"1"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Role     string `json:"role" example:"user"`
}

// AdminResponse is returned by the admin check endpoint.
type AdminResponse struct {
	Message string `json:"message" example:"welcome admin"`
}

// HealthResponse is returned by the health check endpoint.
type HealthResponse struct {
	Status string `json:"status" example:"running"`
}
