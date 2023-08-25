package model

import "time"

type User struct {
	ID        int       `json:"id"`
	CreatedAT time.Time `json:"created_at"`
}

type Segment struct {
	ID          int       `json:"id"`
	SegmentName string    `json:"segment_name"`
	CreatedAT   time.Time `json:"created_at"`
}

type UserSegment struct {
	UserID     int       `json:"user_id"`
	SegmentID  int       `json:"segment_id"`
	AddAT      time.Time `json:"add_at"`
	DeleteTime time.Time `json:"delete_time"`
}
