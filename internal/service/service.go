package service

import (
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

func (s *Service) GetUserSegment(segmentName string) ([]model.UserSegment, error) {
	result, err := s.Storage.GetUserSegment(segmentName)
	return result, err
}

func (s *Service) CreateUserSegment(segment *model.UserSegment) error {
	err := s.Storage.CreateUserSegment(segment)
	if err != nil {
		return err
	}
	return nil
}
