package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {
	mongoDBURL := "mongodb://%s:%s@%s:%s"

	clientOptions := options.Client().ApplyURI(mongoDBURL).SetAuth(options.Credential{
		AuthSource: authDB,
		Username:   username,
		Password:   password,
	})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping to mongodb %v", err)
	}

	return client.Database(database), nil

}
