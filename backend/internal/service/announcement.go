package service

import "time"

// Announcement represents a system announcement.
type Announcement struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	Version     string     `json:"version"`
	Category    string     `json:"category"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
