package user

import (
	"context"
	"rest-api/internal/storage"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage struct {
	collection *mongo.Collection
	logger     *logrus.Logger
}

func NewMongoStorage(client *mongo.Client, dbName, collectionName string, logger *logrus.Logger) *MongoStorage {
	logger.Infof("Initializing MongoStorage for database: %s, collection: %s", dbName, collectionName)
	return &MongoStorage{
		collection: client.Database(dbName).Collection(collectionName),
		logger:     logger,
	}
}

func (s *MongoStorage) GetAll(ctx context.Context) ([]storage.Client, error) {
	s.logger.Info("Fetching all users from the database")

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.logger.Errorf("Failed to fetch users: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []storage.Client
	for cursor.Next(ctx) {
		var user storage.Client
		if err := cursor.Decode(&user); err != nil {
			s.logger.Errorf("Failed to decode user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Errorf("Cursor error while fetching users: %v", err)
		return nil, err
	}

	s.logger.Infof("Successfully fetched %d users", len(users))
	return users, nil
}

func (s *MongoStorage) Create(ctx context.Context, client storage.Client) (string, error) {
	s.logger.Infof("Creating a new user: %+v", client)

	res, err := s.collection.InsertOne(ctx, client)
	if err != nil {
		s.logger.Errorf("Failed to insert user: %v", err)
		return "", err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		s.logger.Error("Inserted ID is not a valid ObjectID")
		return "", mongo.ErrNilDocument
	}

	s.logger.Infof("User created successfully with ID: %s", id.Hex())
	return id.Hex(), nil
}

func (s *MongoStorage) FindOne(ctx context.Context, id string) (storage.Client, error) {
	s.logger.Infof("Fetching user with ID: %s", id)

	var user storage.Client
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("Invalid ObjectID format: %v", err)
		return user, err
	}

	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Warnf("User with ID %s not found", id)
		} else {
			s.logger.Errorf("Failed to fetch user: %v", err)
		}
		return user, err
	}

	s.logger.Infof("User found: %+v", user)
	return user, nil
}

func (s *MongoStorage) Update(ctx context.Context, client storage.Client) error {
	s.logger.Infof("Updating user with ID: %s", client.ID)

	objectID, err := primitive.ObjectIDFromHex(client.ID)
	if err != nil {
		s.logger.Errorf("Invalid ObjectID format: %v", err)
		return err
	}

	result, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"email":    client.Email,
			"username": client.Username,
			"password": client.PasswordHash,
		}},
	)
	if err != nil {
		s.logger.Errorf("Failed to update user: %v", err)
		return err
	}

	s.logger.Infof("User updated successfully, modified count: %d", result.ModifiedCount)
	return nil
}

func (s *MongoStorage) PartiallyUpdate(ctx context.Context, client storage.Client) error {
	s.logger.Infof("Partially updating user with ID: %s", client.ID)

	objectID, err := primitive.ObjectIDFromHex(client.ID)
	if err != nil {
		s.logger.Errorf("Invalid ObjectID format: %v", err)
		return err
	}

	updateFields := bson.M{}
	if client.Email != "" {
		updateFields["email"] = client.Email
	}
	if client.Username != "" {
		updateFields["username"] = client.Username
	}
	if client.PasswordHash != "" {
		updateFields["password"] = client.PasswordHash
	}

	result, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		s.logger.Errorf("Failed to partially update user: %v", err)
		return err
	}

	s.logger.Infof("User partially updated successfully, modified count: %d", result.ModifiedCount)
	return nil
}

func (s *MongoStorage) Delete(ctx context.Context, id string) error {
	s.logger.Infof("Deleting user with ID: %s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Errorf("Invalid ObjectID format: %v", err)
		return err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		s.logger.Errorf("Failed to delete user: %v", err)
		return err
	}

	s.logger.Infof("User deleted successfully, deleted count: %d", result.DeletedCount)
	return nil
}
