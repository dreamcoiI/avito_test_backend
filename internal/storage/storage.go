package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
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

	rows, err := s.db.QueryContext(ctx, query, id)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var segmentName string
		err := rows.Scan(&segmentName)
		if err != nil {
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

func (s *Storage) DeleteSegmentToUser(ctx context.Context, delete []string, id int) error {

	var count int
	count, err := s.CheckUser(id)
	if count < 1 {
		return fmt.Errorf("users with id '%d' not found", id)
	}
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, segmentName := range delete {
		exists, err := s.CheckUserSegmentExists(ctx, id, segmentName)
		if err != nil {
			return err
		}

		if exists {
			err := s.DeleteUserSegment(ctx, id, segmentName)
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
	query := `
		INSERT INTO segment_user (id_user, id_segment)
		SELECT $1, id FROM segments WHERE segment_name = $2 AND id NOT IN (
			SELECT id_segment FROM segment_user WHERE id_user = $1 AND delete_time IS NOT NULL
		)
	`
	_, err := s.db.ExecContext(ctx, query, userID, segmentName)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteUserSegment(ctx context.Context, userID int, segmentName string) error {
	query := `
		UPDATE segment_user
		SET delete_time = CURRENT_TIMESTAMP
		WHERE id_user = $1 AND id_segment = (SELECT id FROM segments WHERE segment_name = $2) AND delete_time is NULL
	`

	_, err := s.db.ExecContext(ctx, query, userID, segmentName)
	if err != nil {
		return fmt.Errorf("DeleteUserSegment: failed to execute query: %w", err)
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
func (s *Storage) GenerateSegmentHistoryCSVByMonth(ctx context.Context, year, month int, filename string) (string, error) {
	startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	query := `
			SELECT id_user, id_segment, 
		   CASE 
			   WHEN delete_time IS NOT NULL THEN 'удаление'
			   ELSE 'добавление'
		   END AS операция,
		   coalesce(delete_time, add_at) AS дата_и_время
	FROM segment_user
	WHERE (add_at >= $1 AND add_at <= $2) OR (delete_time >= $1 AND delete_time <= $2)
		`

	rows, err := s.db.QueryContext(ctx, query, startTime, endTime)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var csvData string
	csvData += "идентификатор пользователя;сегмент;операция;дата и время\n"

	for rows.Next() {
		var userID int
		var segmentName, operation string
		var timestamp time.Time
		err := rows.Scan(&userID, &segmentName, &operation, &timestamp)
		if err != nil {
			return "", err
		}

		csvData += fmt.Sprintf("%d;%s;%s;%s\n", userID, segmentName, operation, timestamp.Format("2006-01-02 15:04:05"))
	}

	filePath := "" + filename
	err = saveCSVToFile(csvData, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func saveCSVToFile(csvData string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(csvData)
	if err != nil {
		return err
	}

	return nil
}
