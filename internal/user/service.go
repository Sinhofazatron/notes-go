package user

import (
	"context"
	"notes-go/pkg/logging"
)

type Service struct {
	storage Storage
	logger *logging.Logger
}

func (s *Service) Create(ctx context.Context, dto  CreateUserDTO) (u User, err error) {
	// TODO: for next one
	return
}