package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"todo-app/internal/auth"
	"todo-app/internal/config"
	"todo-app/internal/controller"
	"todo-app/internal/repository"
	"todo-app/internal/routes"
	"todo-app/internal/service"
	"todo-app/pkg/database"

	_ "todo-app/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Todo API
// @version 1.0
// @description This is a simple todo app backend API.

// @contact.name Muhammed Ibrahim
// @contact.email aminmuhammad18@gmail.com

// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Initialize repositories
	todoRepo := repository.NewTodoRepository(mongoDB.Database, "todos")
	userRepo := repository.NewUserRepository(mongoDB.Database, "users")

	// Initialize services
	authService := auth.NewAuthService(cfg.JWTSecret, cfg.JWTExpiration, cfg.PasswordPepper, userRepo)
	todoService := service.NewTodoService(todoRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	todoController := controller.NewTodoController(todoService)

	// Set up Gin
	if cfg.TestMode {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Set up routes
	routes.SetupRoutes(router, authController, todoController, authService)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	go func() {
		if err := router.Run(cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
