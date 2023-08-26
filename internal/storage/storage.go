package storage

import (
	"context"
	"fmt"
	"github.com/dreamcoiI/avito_test_backend/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DataBasePool *pgxpool.Pool
}

func NewStorage(dataBase *pgxpool.Pool) *Storage {
	storage := new(Storage)
	storage.DataBasePool = dataBase
	return storage
}

func (database *Storage) FindUserSegment(segmentName string) ([]model.UserSegment, error) {
	query := "SELECT id_user, id_segment, created_at FROM users where name = $1"

	var result []model.UserSegment

	err := pgxscan.Select(context.Background(), database.DataBasePool, &result, query, segmentName)
	if err != nil {
		return nil, fmt.Errorf("FindUserSegment: Error %w ", err)
	}
	return result, err
}
