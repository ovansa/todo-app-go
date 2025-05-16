package repository

import (
	"context"
	"errors"
	"time"
	"todo-app/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *model.TodoCreate) (*model.Todo, error)
	FindByID(ctx context.Context, id string, userId string) (*model.Todo, error)
	FindAll(ctx context.Context, userId string) ([]*model.Todo, error)
	Update(ctx context.Context, id string, userId string, todo *model.TodoUpdate) (*model.Todo, error)
	Delete(ctx context.Context, id string, userId string) error
}

type todoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(db *mongo.Database, collectionName string) TodoRepository {
	return &todoRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *todoRepository) Create(ctx context.Context, todoCreate *model.TodoCreate) (*model.Todo, error) {
	userID, ok := ctx.Value("userId").(primitive.ObjectID)
	if !ok {
		return nil, errors.New("user ID not found in context or invalid format")
	}

	now := time.Now()

	completed := false
	if todoCreate.Completed != nil {
		completed = *todoCreate.Completed
	}

	todo := &model.Todo{
		Title:     todoCreate.Title,
		Completed: completed,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
	}

	result, err := r.collection.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)
	return todo, nil
}

func (r *todoRepository) FindAll(ctx context.Context, userId string) ([]*model.Todo, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": userObjectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []*model.Todo
	if err := cursor.All(ctx, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *todoRepository) FindByID(ctx context.Context, id string, userId string) (*model.Todo, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	var todo model.Todo
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID, "userId": userObjectID}).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	return &todo, nil
}

func (r *todoRepository) Update(ctx context.Context, id string, userId string, updateData *model.TodoUpdate) (*model.Todo, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	updateData.UpdatedAt = time.Now()
	update := bson.M{"$set": updateData}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID, "userId": userObjectID}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("todo not found")
	}

	return r.FindByID(ctx, id, userId)
}

func (r *todoRepository) Delete(ctx context.Context, id string, userId string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return errors.New("invalid user id format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID, "userId": userObjectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("todo not found")
	}

	return nil
}
