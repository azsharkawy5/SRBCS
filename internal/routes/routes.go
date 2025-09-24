package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/azsharkawy5/SRBCS/internal/handler"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(engine *gin.Engine, userHandler *handler.UserHandler) {
	// API version prefix
	api := engine.Group("/api/v1")

	// Health check endpoint (no auth required)
	api.GET("/health", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.Writer.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// User routes
	users := api.Group("/users")
	{
		users.POST("/", userHandler.CreateUser)
		users.GET("/", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	// Debug routes (in development only)
	if os.Getenv("APP_ENV") == "development" {
		debug := api.Group("/debug")
		debug.GET("/routes", func(c *gin.Context) {
			c.Header("Content-Type", "application/json")
			c.Status(200)
			response := `{"routes": ["GET /api/v1/health", "POST /api/v1/users", "GET /api/v1/users", "GET /api/v1/users/:id", "PUT /api/v1/users/:id", "DELETE /api/v1/users/:id"]}`
			c.Writer.Write([]byte(response))
		})
	}
}
