package validators

import (
	"net/mail"
	"strings"
)

// SignupRequest represents the request body for user signup.
type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ValidateSignup validates the signup request.
func ValidateSignup(req SignupRequest) []ValidationError {
	var errors []ValidationError

	// Name is required
	name := strings.TrimSpace(req.Name)
	if name == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	} else if len(name) < 2 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must be at least 2 characters",
		})
	} else if len(name) > 100 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must be at most 100 characters",
		})
	}

	// Email is required and must be valid
	email := strings.TrimSpace(req.Email)
	if email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	} else if _, err := mail.ParseAddress(email); err != nil {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Email must be a valid email address",
		})
	}

	// Password is required with minimum length
	if req.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	} else if len(req.Password) < 6 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must be at least 6 characters",
		})
	} else if len(req.Password) > 128 {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password must be at most 128 characters",
		})
	}

	return errors
}

// ValidateLogin validates the login request.
func ValidateLogin(req LoginRequest) []ValidationError {
	var errors []ValidationError

	email := strings.TrimSpace(req.Email)
	if email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	}

	if req.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Message: "Password is required",
		})
	}

	return errors
}
