package service

import (
	"context"
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

func (s *Service) CreateSegment(ctx context.Context, slug string) error {
	err := s.Storage.CreateSegments(ctx, slug)
	return err
}

func (s *Service) DeleteSegment(ctx context.Context, slug string) error {
	err := s.Storage.DeleteSegment(ctx, slug)
	return err
}

func (s *Service) AddSegmentToUser(ctx context.Context, adds []string, id int) error {
	err := s.Storage.AddSegmentToUser(ctx, adds, id)
	if err != nil {
		return err
	}
	return nil
}
