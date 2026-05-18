package auth

import (
	"errors"
	"regexp"
	"unicode"
)

var usernameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_.]{2,31}$`)

func ValidateUsername(username string) error {
	if !usernameRe.MatchString(username) {
		return errors.New("username must be 3-32 chars and use only a-z, 0-9, '_' or '.'")
	}
	if contains(username, "..") || contains(username, "__") || contains(username, "._") || contains(username, "_.") {
		return errors.New("username contains invalid sequence")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters")
	}
	if len(password) > 72 {
		return errors.New("password must be at most 72 characters")
	}

	var hasLower, hasUpper, hasDigit, hasSymbol bool
	for _, r := range password {
		if unicode.IsSpace(r) {
			return errors.New("password must not contain spaces")
		}
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasSymbol = true
		}
	}

	if !(hasLower && hasUpper && hasDigit && hasSymbol) {
		return errors.New("password must include lowercase, uppercase, number, and symbol")
	}
	return nil
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
