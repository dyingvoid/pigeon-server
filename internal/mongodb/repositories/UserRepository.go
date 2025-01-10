package repositories

import (
	"context"
	"fmt"
	"github.com/dyingvoid/pigeon-server/internal/application/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func (r *UserRepository) Add(ctx context.Context, user models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByName(ctx context.Context, name string) (models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0)
	filter := bson.M{}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		err = cursor.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("failed ot decode user: %w", err)
		}
		users = append(users, user)
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		collection: collection,
	}
}
