package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/azsharkawy5/SRBCS/internal/domain"
	"github.com/azsharkawy5/SRBCS/internal/repository/dto"
)

// PostgresUserRepository implements the UserRepository interface
type PostgresUserRepository struct {
	db *sqlx.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// Create inserts a new user into the database and returns the generated ID
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	// Convert domain user to DTO
	userDTO := dto.FromDomain(user)

	query := `
		INSERT INTO users (email, name, is_email_verified, is_active, otp, otp_expires_at, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var generatedID string
	err := r.db.QueryRowContext(ctx, query,
		userDTO.Email,
		userDTO.Name,
		userDTO.IsEmailVerified,
		userDTO.IsActive,
		userDTO.OTP,
		userDTO.OTPExpiresAt,
		userDTO.Role,
		userDTO.CreatedAt,
		userDTO.UpdatedAt,
	).Scan(&generatedID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Set the generated ID back to the domain user object
	user.ID = generatedID
	return nil
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, name, is_email_verified, is_active, otp, otp_expires_at, role, created_at, updated_at
		FROM users
		WHERE id = $1`

	var userDTO dto.UserDTO
	err := r.db.GetContext(ctx, &userDTO, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Convert DTO to domain user
	user, err := userDTO.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert user DTO to domain: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, name, is_email_verified, is_active, otp, otp_expires_at, role, created_at, updated_at
		FROM users
		WHERE email = $1`

	var userDTO dto.UserDTO
	err := r.db.GetContext(ctx, &userDTO, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Convert DTO to domain user
	user, err := userDTO.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert user DTO to domain: %w", err)
	}

	return user, nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	// Convert domain user to DTO
	userDTO := dto.FromDomain(user)

	query := `
		UPDATE users
		SET email = $2, name = $3, is_email_verified = $4, is_active = $5, otp = $6, otp_expires_at = $7, role = $8, updated_at = $9
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query,
		userDTO.ID,
		userDTO.Email,
		userDTO.Name,
		userDTO.IsEmailVerified,
		userDTO.IsActive,
		userDTO.OTP,
		userDTO.OTPExpiresAt,
		userDTO.Role,
		userDTO.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete removes a user from the database
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List retrieves a paginated list of users as DTOs
func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, email, name, is_email_verified, is_active, otp, otp_expires_at, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	var usersDTO []dto.UserDTO
	err := r.db.SelectContext(ctx, &usersDTO, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	var users []*domain.User
	for _, userDTO := range usersDTO {
		user, err := userDTO.ToDomain()
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
