package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		userName string
		wantErr  bool
		errType  error
	}{
		{
			name:     "valid user",
			email:    "test@example.com",
			userName: "Test User",
			wantErr:  false,
		},
		{
			name:     "empty name",
			email:    "test@example.com",
			userName: "",
			wantErr:  true,
			errType:  ErrInvalidUserName,
		},
		{
			name:     "invalid email format",
			email:    "invalid-email",
			userName: "Test User",
			wantErr:  true,
			errType:  ErrInvalidUserEmail,
		},
		{
			name:     "empty email",
			email:    "",
			userName: "Test User",
			wantErr:  true,
			errType:  ErrInvalidUserEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.userName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error, got nil")
					return
				}

				// Check for specific error type if provided
				if tt.errType != nil {
					// Since we're wrapping errors, we need to check if the target error is contained
					if !containsTargetError(err, tt.errType) {
						t.Errorf("NewUser() expected error %v, got %v", tt.errType, err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("NewUser() returned nil user")
				return
			}

			// Validate user fields
			if user.Email != tt.email {
				t.Errorf("NewUser() Email = %v, want %v", user.Email, tt.email)
			}

			if user.Name != tt.userName {
				t.Errorf("NewUser() Name = %v, want %v", user.Name, tt.userName)
			}

			// Check timestamps
			if user.CreatedAt.IsZero() {
				t.Errorf("NewUser() CreatedAt should not be zero")
			}

			if user.UpdatedAt.IsZero() {
				t.Errorf("NewUser() UpdatedAt should not be zero")
			}

			// CreatedAt and UpdatedAt should be close in time
			timeDiff := user.UpdatedAt.Sub(user.CreatedAt)
			if timeDiff > time.Second {
				t.Errorf("NewUser() CreatedAt and UpdatedAt differ by %v, expected less than 1 second", timeDiff)
			}
		})
	}
}

func TestUser_IsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "valid email",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "valid email with subdomain",
			email: "user@mail.example.com",
			want:  true,
		},
		{
			name:  "valid email with numbers",
			email: "user123@example123.com",
			want:  true,
		},
		{
			name:  "valid email with plus",
			email: "user+test@example.com",
			want:  true,
		},
		{
			name:  "invalid email - no @",
			email: "userexample.com",
			want:  false,
		},
		{
			name:  "invalid email - no domain",
			email: "user@",
			want:  false,
		},
		{
			name:  "invalid email - no local part",
			email: "@example.com",
			want:  false,
		},
		{
			name:  "invalid email - no TLD",
			email: "user@example",
			want:  false,
		},
		{
			name:  "empty email",
			email: "",
			want:  false,
		},
		{
			name:  "email with spaces",
			email: "user @example.com",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Email: tt.email}
			if got := user.IsValidEmail(); got != tt.want {
				t.Errorf("User.IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	user := &User{
		ID:        "user-1",
		Email:     "old@example.com",
		Name:      "Test User",
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	originalUpdatedAt := user.UpdatedAt

	tests := []struct {
		name     string
		newEmail string
		wantErr  bool
		errType  error
	}{
		{
			name:     "valid email update",
			newEmail: "new@example.com",
			wantErr:  false,
		},
		{
			name:     "invalid email update",
			newEmail: "invalid-email",
			wantErr:  true,
			errType:  ErrInvalidUserEmail,
		},
		{
			name:     "empty email update",
			newEmail: "",
			wantErr:  true,
			errType:  ErrInvalidUserEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset user email for each test
			user.Email = "old@example.com"
			user.UpdatedAt = originalUpdatedAt

			err := user.UpdateEmail(tt.newEmail)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateEmail() expected error, got nil")
					return
				}

				if tt.errType != nil && err != tt.errType {
					t.Errorf("UpdateEmail() expected error %v, got %v", tt.errType, err)
				}

				// Email should remain unchanged on error
				if user.Email != "old@example.com" {
					t.Errorf("UpdateEmail() email should remain unchanged on error, got %v", user.Email)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateEmail() unexpected error: %v", err)
				return
			}

			if user.Email != tt.newEmail {
				t.Errorf("UpdateEmail() Email = %v, want %v", user.Email, tt.newEmail)
			}

			// UpdatedAt should be changed
			if !user.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("UpdateEmail() UpdatedAt should be updated")
			}
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	user := &User{
		Email:     "test@example.com",
		Name:      "Old Name",
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	originalUpdatedAt := user.UpdatedAt

	tests := []struct {
		name    string
		newName string
		wantErr bool
		errType error
	}{
		{
			name:    "valid name update",
			newName: "New Name",
			wantErr: false,
		},
		{
			name:    "empty name update",
			newName: "",
			wantErr: true,
			errType: ErrInvalidUserName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset user name for each test
			user.Name = "Old Name"
			user.UpdatedAt = originalUpdatedAt

			err := user.UpdateName(tt.newName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateName() expected error, got nil")
					return
				}

				if tt.errType != nil && err != tt.errType {
					t.Errorf("UpdateName() expected error %v, got %v", tt.errType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateName() unexpected error: %v", err)
				return
			}

			if user.Name != tt.newName {
				t.Errorf("UpdateName() Name = %v, want %v", user.Name, tt.newName)
			}

			// UpdatedAt should be changed
			if !user.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("UpdateName() UpdatedAt should be updated")
			}
		})
	}
}

// Helper function to check if an error contains a specific target error
func containsTargetError(err, target error) bool {
	return errors.Is(err, target)
}
