package domain

import "errors"

// User-related errors
var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrInvalidUserID            = errors.New("invalid user ID")
	ErrInvalidUserEmail         = errors.New("invalid user email")
	ErrInvalidUserName          = errors.New("invalid user name")
	ErrInvalidUserRole          = errors.New("invalid user role")
	ErrInvalidUserEmailVerified = errors.New("invalid user email verified")
	ErrInvalidUserActive        = errors.New("invalid user active")
	ErrInvalidOTP               = errors.New("invalid OTP")
	ErrInvalidOTPExpiresAt      = errors.New("OTP expires at is in the past")
)

var (
	ErrInternalError    = errors.New("internal server error")
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrValidationFailed = errors.New("validation failed")
)
