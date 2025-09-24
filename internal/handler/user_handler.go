package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/azsharkawy5/SRBCS/internal/domain"
	"github.com/azsharkawy5/SRBCS/internal/repository/dto"
)

// UserService interface defines what the handler needs from the service layer
type UserService interface {
	CreateUser(ctx context.Context, email, name string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, email, name string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// UserResponse represents the response body for user operations
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.writeError(c, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	// Validate required fields
	if req.Email == "" || req.Name == "" {
		h.writeError(c, http.StatusBadRequest, "Missing required fields", "email and name are required")
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Email, req.Name)
	if err != nil {
		statusCode := h.getStatusCodeFromError(err)
		h.writeError(c, statusCode, "Failed to create user", err.Error())
		return
	}

	response := h.userToResponse(user)
	c.JSON(http.StatusCreated, response)
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.writeError(c, http.StatusBadRequest, "Missing user ID", "")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		statusCode := h.getStatusCodeFromError(err)
		h.writeError(c, statusCode, "Failed to get user", err.Error())
		return
	}

	response := h.userToResponse(user)
	c.JSON(http.StatusOK, response)
}

// UpdateUser handles PUT /users/{id}
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.writeError(c, http.StatusBadRequest, "Missing user ID", "")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.writeError(c, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req.Email, req.Name)
	if err != nil {
		statusCode := h.getStatusCodeFromError(err)
		h.writeError(c, statusCode, "Failed to update user", err.Error())
		return
	}

	response := h.userToResponse(user)
	c.JSON(http.StatusOK, response)
}

// DeleteUser handles DELETE /users/{id}
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.writeError(c, http.StatusBadRequest, "Missing user ID", "")
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		statusCode := h.getStatusCodeFromError(err)
		h.writeError(c, statusCode, "Failed to delete user", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse query parameters
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 10 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default offset
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	users, err := h.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		statusCode := h.getStatusCodeFromError(err)
		h.writeError(c, statusCode, "Failed to list users", err.Error())
		return
	}

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = h.userToResponse(user)
	}

	c.JSON(http.StatusOK, responses)
}

// userToResponse converts a domain user to response format
func (h *UserHandler) userToResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// dtoToResponse converts a user DTO to response format
func (h *UserHandler) dtoToResponse(userDTO *dto.UserDTO) UserResponse {
	return UserResponse{
		ID:        userDTO.ID,
		Email:     userDTO.Email,
		Name:      userDTO.Name,
		CreatedAt: userDTO.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: userDTO.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// getStatusCodeFromError maps domain errors to HTTP status codes
func (h *UserHandler) getStatusCodeFromError(err error) int {
	switch {
	case containsError(err, domain.ErrUserNotFound):
		return http.StatusNotFound
	case containsError(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict
	case containsError(err, domain.ErrInvalidUserID),
		containsError(err, domain.ErrInvalidUserEmail),
		containsError(err, domain.ErrInvalidUserName),
		containsError(err, domain.ErrInvalidInput),
		containsError(err, domain.ErrValidationFailed):
		return http.StatusBadRequest
	case containsError(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized
	case containsError(err, domain.ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// writeError writes an error response
func (h *UserHandler) writeError(c *gin.Context, statusCode int, errTitle, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error:   errTitle,
		Message: message,
	})
}
