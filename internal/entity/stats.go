package entity

import "time"

type Stats struct {
	Id                     int       `db:"id" json:"id"`
	ProcessedCount         int       `db:"processed_count" json:"processed_count"`
	LastProcessedContent   string    `db:"last_processed_message_content" json:"last_processed_content"`
	LastProcessedMessageId int       `db:"last_processed_message_id" json:"last_processed_message_id"`
	UpdatedAt              time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
