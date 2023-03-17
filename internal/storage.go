package internal

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, task Task) (string, error)
	FindOne(ctx context.Context, id string) (Task, error)
	FindAll(ctx context.Context) (t []Task, err error)
	UpdateTask(ctx context.Context, task Task) error
	DeleteTask(ctx context.Context, id string) error
}
