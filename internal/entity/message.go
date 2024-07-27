package entity

import "time"

type Message struct {
	Id          int       `db:"id"`
	Content     string    `db:"content"`
	IsProcessed bool      `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	ProcessedAt time.Time `db:"processed_at"`
}
