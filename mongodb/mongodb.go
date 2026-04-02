package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDB struct {
	client *mongo.Client
}

func (m *MongoDB) Save(originalURL, shortenURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.client.Database("urlshortener").Collection("urls")
	_, err := collection.InsertOne(ctx, bson.M{"key": shortenURL, "original_url": originalURL})

	return err
}

func (m *MongoDB) OriginalURL(shortenURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.client.Database("urlshortener").Collection("urls")
	var result result
	err := collection.FindOne(ctx, bson.M{"key": shortenURL}).Decode(&result)
	return result.OriginalURL, err
}

type result struct {
	OriginalURL string `bson:"original_url"`
}
