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
	GetAll(ctx context.Context) ([]User, error)
	PartiallyUpdate(ctx context.Context, user User) error
}

type MongoStorage struct {
	collection *mongo.Collection
}

func NewMongoStorage(client *mongo.Client, dbName, collectionName string) *MongoStorage {
	return &MongoStorage{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (s *MongoStorage) GetAll(ctx context.Context) ([]User, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoStorage) Create(ctx context.Context, user User) (string, error) {
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", mongo.ErrNilDocument
	}

	return id.Hex(), nil
}

func (s *MongoStorage) FindOne(ctx context.Context, id string) (User, error) {
	var user User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *MongoStorage) Update(ctx context.Context, user User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"email":    user.Email,
			"username": user.Username,
			"password": user.PasswordHash,
		}},
	)
	return err
}

func (s *MongoStorage) PartiallyUpdate(ctx context.Context, user User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	updateFields := bson.M{}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	if user.Username != "" {
		updateFields["username"] = user.Username
	}
	if user.PasswordHash != "" {
		updateFields["password"] = user.PasswordHash
	}

	_, err = s.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	return err
}

func (s *MongoStorage) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
