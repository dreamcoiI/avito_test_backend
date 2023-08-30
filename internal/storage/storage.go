package storage

import (
	"context"
	"database/sql"
	"fmt"
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

	var count int
	count, err := s.CheckSegment(slug)
	if count > 0 {
		return fmt.Errorf("segment with name '%s' already exists", slug)
	}

	query := "INSERT INTO segments VALUES ((SELECT max(id) +1 from segments), $1) "

	fmt.Println("Denis aboba")

	_, err = s.db.ExecContext(ctx, query, slug)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *Storage) DeleteSegment(ctx context.Context, slug string) error {

	var count int
	count, err := s.CheckSegment(slug)
	if count == 0 {
		return fmt.Errorf("segment with name '%s' not found", slug)
	}

	query := "DELETE FROM segments WHERE segment_name = $1"
	_, err = s.db.ExecContext(ctx, query, slug)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}

func (s *Storage) AddSegmentToUser(ctx context.Context, adds []string, id int) error {

	var count int
	count, err := s.CheckUser(id)
	if count < 1 {
		return fmt.Errorf("users with id '%d' not found", id)
	}
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, segmentName := range adds {
		exists, err := s.CheckUserSegmentExists(ctx, id, segmentName)
		if err != nil {
			return err
		}
		if !exists {
			err := s.InsertUserSegment(ctx, id, segmentName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Storage) CheckUserSegmentExists(ctx context.Context, userID int, segmentName string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM segment_user su JOIN segments s ON su.id_segment = s.id WHERE su.id_user = $1 AND s.segment_name = $2)"

	var exists bool
	err := s.db.QueryRowContext(ctx, query, userID, segmentName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) InsertUserSegment(ctx context.Context, userID int, segmentName string) error {
	query := "INSERT INTO segment_user (id_user, id_segment) SELECT $1, id FROM segments WHERE segment_name = $2"
	_, err := s.db.ExecContext(ctx, query, userID, segmentName)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CheckSegment(slug string) (int, error) {
	existsQuery := "SELECT COUNT(*) FROM segments WHERE segment_name = $1"
	var count int
	err := s.db.QueryRow(existsQuery, slug).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, err
}

func (s *Storage) CheckUser(id int) (int, error) {
	existsQuery := "SELECT COUNT(*) FROM users WHERE id = $1"
	var count int
	err := s.db.QueryRow(existsQuery, id).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, err
}
