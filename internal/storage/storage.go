package storage

import (
	"github.com/dreamcoiI/avito_test_backend/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	storagePool *pgxpool.Pool
}

func NewStorage(dataBase *pgxpool.Pool) *Storage {
	storage := new(Storage)
	storage.storagePool = dataBase
	return storage
}

func (database *Storage) GetUserSegment(segmentName string) ([]model.UserSegment, error) {
	//query := "SELECT segments.id, segments.segment_name,segments.created_at FROM segments JOIN segment_user ON segment_ust.id_segment = segments.id  WHERE id_user = $1" // TODO right query
	//
	var result []model.UserSegment
	//
	//err := pgxscan.Select(context.Background(), database.storagePool, &result, query, segmentName)
	//if err != nil {
	//	return nil, fmt.Errorf("FindUserSegment: Error %w ", err)
	//}
	return result, nil
}

func (database *Storage) CreateUserSegment(segment *model.UserSegment) error {
	//
	return nil
}
