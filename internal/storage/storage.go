package storage

import "context"

type Storage interface {
	Create(ctx context.Context, client Client) (string, error)
	FindOne(ctx context.Context, id string) (Client, error)
	Update(ctx context.Context, client Client) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]Client, error)
	PartiallyUpdate(ctx context.Context, client Client) error
}
