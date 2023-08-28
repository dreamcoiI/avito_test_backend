package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dreamcoiI/avito_test_backend/internal/model"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(conn *sql.DB) *Storage {
	storage := &Storage{}
	storage.db = conn
	return storage
}

func (s *Storage) GetUserSegment(ctx context.Context, id int) (string, error) {
	query := "SELECT segments.segment_name FROM segments JOIN segment_user ON segment_user.id_segment = segments.id WHERE id_user = $1"

	var segmentName string
	//time.Sleep(5 * time.Second)
	fmt.Println(id)
	err := s.db.QueryRowContext(ctx, query, int64(id)).Scan(&segmentName)
	if err == sql.ErrNoRows {
		fmt.Println("Денис Абоба")
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed on db.QueryRowContext: %w ", err)
	}
	fmt.Println(segmentName)
	return segmentName, nil
}

func (s *Storage) CreateUserSegment(segment *model.UserSegment) error {
	//
	return nil
}
