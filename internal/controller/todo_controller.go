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

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param todo body model.TodoCreate true "Todo details"
// @Success 201 {object} model.Todo
// @Router /todos [post]
func (c *TodoController) CreateTodo(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in context"})
		return
	}

	var todoCreate model.TodoCreate
	if err := ctx.ShouldBindJSON(&todoCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTodo, err := c.service.CreateTodo(ctx.Request.Context(), userId.(string), &todoCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdTodo)
}

// GetTodo godoc
// @Summary Get a single todo
// @Description Get a todo item by ID
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} model.Todo
// @Router /todos/{id} [get]
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

// GetAllTodos godoc
// @Summary Get all todos
// @Description Retrieve all todos for the authenticated user
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Todo
// @Router /todos [get]
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

// UpdateTodo godoc
// @Summary Update a todo
// @Description Update a todo item by ID
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Param todo body model.TodoUpdate true "Updated todo data"
// @Success 200 {object} model.Todo
// @Router /todos/{id} [put]
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

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Delete a todo item by ID
// @Tags todos
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 204 {object} nil
// @Router /todos/{id} [delete]
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
