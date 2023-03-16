package internal

import (
	"context"
)

type Storage interface {
	CreateTask(ctx context.Context, task ToDoList) (string, error)
	FindOne(ctx context.Context, id string) (ToDoList, error)
	UpdateTask(ctx context.Context, task CreateListDTO) error
	DeleteTask(ctx context.Context, id string) error
}
