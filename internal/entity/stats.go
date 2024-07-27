package entity

import "time"

type Stats struct {
	Id                     int       `db:"id"`
	ProcessedCount         int       `db:"processed_count"`
	LastProcessedMessageId int       `db:"last_processed_message_id"`
	UpdatedAt              time.Time `db:"updated_at"`
}
