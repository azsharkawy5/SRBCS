package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azsharkawy5/SRBCS/config"
	"github.com/azsharkawy5/SRBCS/internal/handler"
	"github.com/azsharkawy5/SRBCS/internal/repository"
	"github.com/azsharkawy5/SRBCS/internal/routes"
	"github.com/azsharkawy5/SRBCS/internal/service"
	"github.com/azsharkawy5/SRBCS/pkg/httpserver"
	"github.com/azsharkawy5/SRBCS/pkg/postgres"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	dbConfig := postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	dbConn, err := postgres.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Check database health
	if err := dbConn.Health(); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}
	log.Println("Database connection established successfully")

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(dbConn.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Initialize HTTP server
	serverConfig := httpserver.Config{
		Host:         cfg.Server.Host,
		Port:         cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	server := httpserver.NewServer(serverConfig)
	engine := server.Engine()

	// Register routes
	routes.RegisterRoutes(engine, userHandler)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// ginWrapHTTPMiddleware adapts a net/http middleware (func(http.Handler) http.Handler) to gin.HandlerFunc
func ginWrapHTTPMiddleware(mw func(http.Handler) http.Handler) func(*gin.Context) {
	return func(c *gin.Context) {
		h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// when the wrapped handler completes, continue to next gin handler
			c.Next()
		}))
		h.ServeHTTP(c.Writer, c.Request)
	}
}
