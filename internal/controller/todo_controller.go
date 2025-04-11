package controller

import (
	"net/http"
	"todo-app/internal/model"
	"todo-app/internal/service"

	"github.com/gin-gonic/gin"
)

type TodoController struct {
	service service.TodoService
}

func NewTodoController(service service.TodoService) *TodoController {
	return &TodoController{service: service}
}

func (c *TodoController) CreateTodo(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	var todo model.Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTodo, err := c.service.CreateTodo(ctx.Request.Context(), userId.(string), &todo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdTodo)
}

func (c *TodoController) GetTodo(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	id := ctx.Param("id")

	todo, err := c.service.GetTodo(ctx.Request.Context(), id, userId.(string))
	if err != nil {
		if err.Error() == "todo not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (c *TodoController) GetAllTodos(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	todos, err := c.service.GetAllTodos(ctx.Request.Context(), userId.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todos)
}

func (c *TodoController) UpdateTodo(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	id := ctx.Param("id")

	var updateData model.TodoUpdate
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTodo, err := c.service.UpdateTodo(ctx.Request.Context(), id, userId.(string), &updateData)
	if err != nil {
		if err.Error() == "todo not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedTodo)
}

func (c *TodoController) DeleteTodo(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	id := ctx.Param("id")

	err := c.service.DeleteTodo(ctx.Request.Context(), id, userId.(string))
	if err != nil {
		if err.Error() == "todo not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
