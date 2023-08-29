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

	fmt.Println(segmentNames)
	return segmentNames, nil
}

func (s *Storage) CreateSegments(ctx context.Context, slug string) error {
	query := "INSERT INTO segments VALUES ((SELECT max(id) +1 from segments), $1) "

	fmt.Println("Denis aboba")

	_, err := s.db.ExecContext(ctx, query, slug)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *Storage) CreateUserSegment(segment *model.UserSegment) error {
	//
	return nil
}
