package dto

import (
	"time"

	"github.com/azsharkawy5/SRBCS/internal/domain"
)

// UserDTO represents the data transfer object for user data in the repository layer
// This DTO is specifically designed for database operations and may include
// database-specific fields that are not part of the domain model
type UserDTO struct {
	ID              string     `db:"id"`
	Email           string     `db:"email"`
	Name            string     `db:"name"`
	IsEmailVerified bool       `db:"is_email_verified"`
	IsActive        bool       `db:"is_active"`
	OTP             *string    `db:"otp"`
	OTPExpiresAt    *time.Time `db:"otp_expires_at"`
	Role            string     `db:"role"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

// ToDomain converts UserDTO to domain.User
func (dto *UserDTO) ToDomain() (*domain.User, error) {
	return domain.NewUserWithID(
		dto.ID,
		dto.Email,
		dto.Name,
		domain.Role(dto.Role),
		dto.IsEmailVerified,
		dto.IsActive,
		dto.OTP,
		dto.OTPExpiresAt,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
}

// FromDomain creates UserDTO from domain.User
func FromDomain(user *domain.User) *UserDTO {
	return &UserDTO{
		ID:              user.ID,
		Email:           user.Email,
		Name:            user.Name,
		IsEmailVerified: user.IsEmailVerified,
		IsActive:        user.IsActive,
		OTP:             user.OTP,
		OTPExpiresAt:    user.OTPExpiresAt,
		Role:            string(user.Role),
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}
