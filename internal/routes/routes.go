package routes

import (
	"todo-app/internal/auth"
	"todo-app/internal/controller"
	"todo-app/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authController *controller.AuthController) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
	}
}

func SetupTodoRoutes(router *gin.Engine, todoController *controller.TodoController, authService auth.Service) {
	todoGroup := router.Group("/todos")
	todoGroup.Use(authService.AuthMiddleware())
	{
		todoGroup.GET("", todoController.GetAllTodos)
		todoGroup.POST("", todoController.CreateTodo)
		todoGroup.GET("/:id", todoController.GetTodo)
		todoGroup.PUT("/:id", todoController.UpdateTodo)
		todoGroup.DELETE("/:id", todoController.DeleteTodo)
	}
}

func SetupRoutes(router *gin.Engine, authController *controller.AuthController, todoController *controller.TodoController, authService auth.Service) {
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORS())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	SetupAuthRoutes(router, authController)
	SetupTodoRoutes(router, todoController, authService)
}
