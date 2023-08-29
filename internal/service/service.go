package service

import (
	"context"
	"github.com/dreamcoiI/avito_test_backend/internal/model"
	"github.com/dreamcoiI/avito_test_backend/internal/storage"
)

type Service struct {
	Storage storage.Storage
}

func NewService(storage *storage.Storage) *Service {
	newService := new(Service)
	newService.Storage = *storage
	return newService
}

func (s *Service) GetUserSegment(ctx context.Context, userID int) ([]string, error) {
	result, err := s.Storage.GetUserSegment(ctx, userID)
	return result, err
}

func (s *Service) CreateUserSegment(segment *model.UserSegment) error {
	err := s.Storage.CreateUserSegment(segment)
	if err != nil {
		return err
	}
	return nil
}
