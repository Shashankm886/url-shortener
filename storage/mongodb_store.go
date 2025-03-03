package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDBStore(uri string, dbName string, collectionName string) (*MongoDBStore, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDBStore{client: client, collection: collection}, err
}

func (m *MongoDBStore) Save(shortURL, longURL string, expiryTime time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.collection.InsertOne(ctx, bson.M{
		"short_url":       shortURL,
		"long_url":        longURL,
		"created_at":      time.Now(),
		"expiration_time": expiryTime,
		"usage":           0,
	})
	return err
}

func (m *MongoDBStore) Get(shortURL string) (*bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	err := m.collection.FindOne(ctx, bson.M{"short_url": shortURL}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MongoDBStore) Delete(shortURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.collection.DeleteOne(ctx, bson.M{"short_url": shortURL})
	return err
}

func (m *MongoDBStore) IncrementUsage(shortURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.collection.UpdateOne(
		ctx,
		bson.M{"short_url": shortURL},
		bson.M{"$inc": bson.M{"usage": 1}},
	)
	return err
}
