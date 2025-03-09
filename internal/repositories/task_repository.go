package repositories

import (
	"context"
	"errors"
	"go-clean-arch/internal/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entities.Task) (string, error)
	GetAll(ctx context.Context) ([]entities.Task, error)
	Update(ctx context.Context, id string, task *entities.Task) error
	Delete(ctx context.Context, id string) error
}

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &taskRepository{
		collection: db.Collection("tasks"),
	}
}

func (r *taskRepository) Create(ctx context.Context, task *entities.Task) (string, error) {
	task.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *taskRepository) GetAll(ctx context.Context) ([]entities.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []entities.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, id string, task *entities.Task) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": task})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("Task not found")
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("Task not found")
	}

	return nil
}
