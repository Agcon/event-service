package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoClient struct {
	client *mongo.Client
}

func ConnectMongo(uri string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &MongoClient{client: client}, nil
}

func (mc *MongoClient) GetDatabase(name string) *mongo.Database {
	return mc.client.Database(name)
}

func (mc *MongoClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mc.client.Disconnect(ctx)
}
