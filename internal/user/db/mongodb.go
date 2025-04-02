package client

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(uri string, logger *logrus.Logger) (*mongo.Client, error) {
	logger.Infof("Initializing connection to MongoDB at URI: %s", uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Errorf("Failed to connect to MongoDB: %v", err)
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	logger.Info("Connected to MongoDB, performing connection check...")

	if err := client.Ping(ctx, nil); err != nil {
		logger.Errorf("MongoDB connection check failed: %v", err)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("MongoDB connection verified successfully")
	return client, nil
}
