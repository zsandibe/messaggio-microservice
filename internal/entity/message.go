package entity

import "time"

type Message struct {
	Id          int       `db:"id" json:"id"`
	Content     string    `db:"content" json:"content"`
	IsProcessed bool      `db:"status" json:"is_processed"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at,omitempty"`
}
