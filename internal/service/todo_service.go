package service

import (
	"context"
	"errors"
	"todo-app/internal/model"
	"todo-app/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoService interface {
	CreateTodo(ctx context.Context, userId string, todo *model.Todo) (*model.Todo, error)
	GetTodo(ctx context.Context, id string, userId string) (*model.Todo, error)
	GetAllTodos(ctx context.Context, userId string) ([]*model.Todo, error)
	UpdateTodo(ctx context.Context, id string, userId string, todo *model.TodoUpdate) (*model.Todo, error)
	DeleteTodo(ctx context.Context, id string, userId string) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(ctx context.Context, userId string, todo *model.Todo) (*model.Todo, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}
	todo.UserID = userObjectID
	return s.repo.Create(ctx, todo)
}

func (s *todoService) GetTodo(ctx context.Context, id string, userId string) (*model.Todo, error) {
	return s.repo.FindByID(ctx, id, userId)
}

func (s *todoService) GetAllTodos(ctx context.Context, userId string) ([]*model.Todo, error) {
	return s.repo.FindAll(ctx, userId)
}

func (s *todoService) UpdateTodo(ctx context.Context, id string, userId string, todo *model.TodoUpdate) (*model.Todo, error) {
	return s.repo.Update(ctx, id, userId, todo)
}

func (s *todoService) DeleteTodo(ctx context.Context, id string, userId string) error {
	return s.repo.Delete(ctx, id, userId)
}
