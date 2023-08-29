package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dreamcoiI/avito_test_backend/internal/model"
	"log"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(conn *sql.DB) *Storage {
	storage := &Storage{}
	storage.db = conn
	return storage
}

func (s *Storage) GetUserSegment(ctx context.Context, id int) ([]string, error) {
	query := "SELECT segments.segment_name FROM segments JOIN segment_user ON segment_user.id_segment = segments.id WHERE id_user = $1"

	var segmentNames []string

	//time.Sleep(5 * time.Second)
	fmt.Println(id)
	rows, err := s.db.Query(query, id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var segmentName string
		err := rows.Scan(&segmentName)
		if err != nil {
			fmt.Println("Денис Абоба")
			log.Fatal(err)
		}
		log.Println(id, segmentName)
		segmentNames = append(segmentNames, segmentName)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//err := s.db.QueryRowContext(ctx, query, int64(id)).Scan(&segmentName)
	//if err == sql.ErrNoRows {
	//	fmt.Println("Денис Абоба")
	//	return "", nil
	//} else if err != nil {
	//	return "", fmt.Errorf("failed on db.QueryRowContext: %w ", err)
	//}
	fmt.Println(segmentNames)
	return segmentNames, nil
} //TODO change GET request for POST ( JSON-body post for bd request on id user for user_id)

func (s *Storage) CreateUserSegment(segment *model.UserSegment) error {
	//
	return nil
}
