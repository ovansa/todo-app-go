package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" binding:"required"`
	Completed bool               `json:"completed" bson:"completed"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type TodoUpdate struct {
	Title       string    `json:"title" bson:"title,omitempty"`
	Description string    `json:"description" bson:"description,omitempty"`
	Completed   bool      `json:"completed" bson:"completed,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}
