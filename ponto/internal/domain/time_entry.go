package domain

import "time"

type EntryType string

const (
	ClockIn  EntryType = "clock_in"
	ClockOut EntryType = "clock_out"
)

type TimeEntry struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Type       EntryType `json:"type"`
	RecordedAt time.Time `json:"recorded_at"`
	Note       string    `json:"note,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
