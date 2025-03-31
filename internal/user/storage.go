package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindOne(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}

type MongoStorage struct {
	collection *mongo.Collection
}

func NewMongoStorage(client *mongo.Client, dbName, collectionName string) *MongoStorage {
	return &MongoStorage{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (s *MongoStorage) Create(ctx context.Context, user User) (string, error) {
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	// Преобразуем InsertedID в строку
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrNilDocument
	}

	return id.Hex(), nil // Преобразуем ObjectID в строку
}

func (s *MongoStorage) FindOne(ctx context.Context, id string) (User, error) {
	var user User
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *MongoStorage) Update(ctx context.Context, user User) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

func (s *MongoStorage) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
