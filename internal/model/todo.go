package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo represents a todo item
// @Description Todo represents a task that a user wants to track
type Todo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"5f8d0614db5c5c7b3a18f201"`
	Title     string             `json:"title" bson:"title" binding:"required" example:"Buy groceries"`
	Completed bool               `json:"completed" bson:"completed" example:"false"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt" example:"2022-01-01T12:00:00Z"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId" example:"5f8d0614db5c5c7b3a18f200"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt" example:"2022-01-01T12:00:00Z"`
}

// TodoCreate is used for creating new todos
// @Description TodoCreate is used when creating a new todo item
type TodoCreate struct {
	Title     string `json:"title" bson:"title" binding:"required" example:"Buy groceries"`
	Completed *bool  `json:"completed" bson:"completed" example:"false"`
}

// TodoUpdate is used for updating existing todos
// @Description TodoUpdate is used when updating an existing todo item
type TodoUpdate struct {
	Title       string    `json:"title" bson:"title,omitempty" example:"Buy more groceries"`
	Description string    `json:"description" bson:"description,omitempty" example:"Need to get milk and eggs"`
	Completed   bool      `json:"completed" bson:"completed,omitempty" example:"true"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt" example:"2022-01-02T12:00:00Z"`
}
