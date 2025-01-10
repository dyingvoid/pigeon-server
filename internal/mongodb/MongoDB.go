package mongodb

import (
	"context"
	"fmt"
	"github.com/dyingvoid/pigeon-server/internal/mongodb/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// TODO context
type Database struct {
	client         *mongo.Client
	UserRepository *repositories.UserRepository
}

type MongoConfig struct {
	ConnectionString   string
	DatabaseName       string
	UserCollectionName string
}

func NewMongoDb(config MongoConfig) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectionString))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	db := client.Database(config.DatabaseName)
	users := db.Collection(config.UserCollectionName)
	err = ensureUserIndexCreated(users)
	if err != nil {
		return nil, fmt.Errorf("could not initialize mongo: %w", err)
	}

	userRepository := repositories.NewUserRepository(users)

	return &Database{
		client:         client,
		UserRepository: userRepository,
	}, nil
}

func ensureUserIndexCreated(collection *mongo.Collection) error {
	publicKeyIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "public_key.der_data", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	nameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(false),
	}

	_, err := collection.Indexes().CreateMany(
		context.Background(),
		[]mongo.IndexModel{
			publicKeyIndex,
			nameIndex,
		},
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}

		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}
