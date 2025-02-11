package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client *mongo.Client
	db     *mongo.Database
}

var client *mongo.Client

func New(storagePath string, dbName string) (*Storage, error) {
	const op = "storage.mongodb.New"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(storagePath))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := client.Database(dbName)

	return &Storage{
		client: client,
		db:     db,
	}, nil
}

func (s *Storage) Database() *mongo.Database {
	return s.db
}
