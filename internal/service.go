package internal

import (
	"app/pkg/logging"
	"context"
)

type service struct {
	storage Storage
	logger  *logging.Logger
}

func (s *service) Create(ctx context.Context, dto CreateListDTO) (task Task, err error) {
	return task, nil
}
