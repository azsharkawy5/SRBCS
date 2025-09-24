package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/azsharkawy5/SRBCS/internal/domain"
)

// MockUserRepository implements UserRepository for testing
type MockUserRepository struct {
	users    map[string]*domain.User
	emails   map[string]*domain.User
	createFn func(ctx context.Context, user *domain.User) error
	getFn    func(ctx context.Context, id string) (*domain.User, error)
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[string]*domain.User),
		emails: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}

	// Check if user with email already exists
	if _, exists := m.emails[user.Email]; exists {
		return errors.New("user already exists")
	}

	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}

	user, exists := m.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, exists := m.emails[email]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	// Update email mapping if email changed
	oldUser := m.users[user.ID]
	if oldUser.Email != user.Email {
		delete(m.emails, oldUser.Email)
		m.emails[user.Email] = user
	}

	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	user, exists := m.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}

	delete(m.users, id)
	delete(m.emails, user.Email)
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}

	// Simple pagination
	start := offset
	if start > len(users) {
		return []*domain.User{}, nil
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		email    string
		userName string
		mockFn   func(*MockUserRepository)
		wantErr  bool
		errType  error
	}{
		{
			name:     "successful user creation",
			email:    "test@example.com",
			userName: "Test User",
			mockFn:   func(m *MockUserRepository) {},
			wantErr:  false,
		},
		{
			name:     "user already exists",
			id:       "user-2",
			email:    "existing@example.com",
			userName: "Existing User",
			mockFn: func(m *MockUserRepository) {
				// Pre-populate with existing user
				existingUser := &domain.User{
					ID:    "existing-user",
					Email: "existing@example.com",
					Name:  "Existing",
				}
				m.users["existing-user"] = existingUser
				m.emails["existing@example.com"] = existingUser
			},
			wantErr: true,
			errType: domain.ErrUserAlreadyExists,
		},
		{
			name:     "invalid email",
			id:       "user-3",
			email:    "invalid-email",
			userName: "Test User",
			mockFn:   func(m *MockUserRepository) {},
			wantErr:  true,
			errType:  domain.ErrInvalidUserEmail,
		},
		{
			name:     "empty name",
			id:       "user-4",
			email:    "test@example.com",
			userName: "",
			mockFn:   func(m *MockUserRepository) {},
			wantErr:  true,
			errType:  domain.ErrInvalidUserName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.mockFn(mockRepo)

			service := NewUserService(mockRepo)

			user, err := service.CreateUser(context.Background(), tt.email, tt.userName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateUser() expected error, got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("CreateUser() expected error %v, got %v", tt.errType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateUser() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("CreateUser() returned nil user")
				return
			}

			if user.ID != tt.id {
				t.Errorf("CreateUser() ID = %v, want %v", user.ID, tt.id)
			}

			if user.Email != tt.email {
				t.Errorf("CreateUser() Email = %v, want %v", user.Email, tt.email)
			}

			if user.Name != tt.userName {
				t.Errorf("CreateUser() Name = %v, want %v", user.Name, tt.userName)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		mockFn  func(*MockUserRepository)
		wantErr bool
		errType error
	}{
		{
			name:   "successful get user",
			userID: "user-1",
			mockFn: func(m *MockUserRepository) {
				user := &domain.User{
					ID:        "user-1",
					Email:     "test@example.com",
					Name:      "Test User",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.users["user-1"] = user
			},
			wantErr: false,
		},
		{
			name:    "user not found",
			userID:  "nonexistent",
			mockFn:  func(m *MockUserRepository) {},
			wantErr: true,
			errType: domain.ErrUserNotFound,
		},
		{
			name:    "empty user ID",
			userID:  "",
			mockFn:  func(m *MockUserRepository) {},
			wantErr: true,
			errType: domain.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			tt.mockFn(mockRepo)

			service := NewUserService(mockRepo)

			user, err := service.GetUserByID(context.Background(), tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserByID() expected error, got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("GetUserByID() expected error %v, got %v", tt.errType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("GetUserByID() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("GetUserByID() returned nil user")
				return
			}

			if user.ID != tt.userID {
				t.Errorf("GetUserByID() ID = %v, want %v", user.ID, tt.userID)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		newEmail string
		newName  string
		mockFn   func(*MockUserRepository) time.Time // Return original UpdatedAt
		wantErr  bool
		errType  error
	}{
		{
			name:     "successful update email and name",
			userID:   "user-1",
			newEmail: "new@example.com",
			newName:  "New Name",
			mockFn: func(m *MockUserRepository) time.Time {
				originalTime := time.Now().Add(-time.Hour)
				existingUser := &domain.User{
					ID:        "user-1",
					Email:     "old@example.com",
					Name:      "Old Name",
					CreatedAt: originalTime,
					UpdatedAt: originalTime,
				}
				m.users["user-1"] = existingUser
				m.emails["old@example.com"] = existingUser
				return originalTime
			},
			wantErr: false,
		},
		{
			name:     "update only name",
			userID:   "user-1",
			newEmail: "",
			newName:  "Updated Name",
			mockFn: func(m *MockUserRepository) time.Time {
				originalTime := time.Now().Add(-time.Hour)
				existingUser := &domain.User{
					ID:        "user-1",
					Email:     "old@example.com",
					Name:      "Old Name",
					CreatedAt: originalTime,
					UpdatedAt: originalTime,
				}
				m.users["user-1"] = existingUser
				m.emails["old@example.com"] = existingUser
				return originalTime
			},
			wantErr: false,
		},
		{
			name:     "user not found",
			userID:   "nonexistent",
			newEmail: "new@example.com",
			newName:  "New Name",
			mockFn:   func(m *MockUserRepository) time.Time { return time.Time{} },
			wantErr:  true,
			errType:  domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			originalUpdatedAt := tt.mockFn(mockRepo)

			service := NewUserService(mockRepo)

			user, err := service.UpdateUser(context.Background(), tt.userID, tt.newEmail, tt.newName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateUser() expected error, got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("UpdateUser() expected error %v, got %v", tt.errType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateUser() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("UpdateUser() returned nil user")
				return
			}

			// Check updated fields
			if tt.newEmail != "" && user.Email != tt.newEmail {
				t.Errorf("UpdateUser() Email = %v, want %v", user.Email, tt.newEmail)
			}

			if tt.newName != "" && user.Name != tt.newName {
				t.Errorf("UpdateUser() Name = %v, want %v", user.Name, tt.newName)
			}

			// Ensure UpdatedAt was changed (compare with original time)
			if !originalUpdatedAt.IsZero() && !user.UpdatedAt.After(originalUpdatedAt) {
				t.Errorf("UpdateUser() UpdatedAt should be newer than original")
			}
		})
	}
}
