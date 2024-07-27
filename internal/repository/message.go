package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zsandibe/messaggio-microservice/internal/domain"
	"github.com/zsandibe/messaggio-microservice/internal/entity"
	logger "github.com/zsandibe/messaggio-microservice/pkg"
)

type messageRepo struct {
	db *sqlx.DB
}

func NewMessageRepo(db *sqlx.DB) *messageRepo {
	return &messageRepo{db: db}
}

func (r *messageRepo) CreateMessage(ctx context.Context, msg domain.CreateMessageRequest) (*entity.Message, error) {
	var message *entity.Message
	var id int
	var status sql.NullBool
	var processed_at sql.NullTime
	startTime := time.Now()

	query := `
	INSERT INTO messages (content) 
	VALUES ($1) 
	RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, msg.Content, status, startTime, processed_at).Scan(&id)
	if err != nil {
		logger.Errorf("Error in inserting message: %v", err)
		return &entity.Message{}, domain.ErrCreatingMessage
	}

	message = &entity.Message{
		Id:          id,
		Content:     msg.Content,
		IsProcessed: status.Bool,
		CreatedAt:   startTime,
		ProcessedAt: processed_at.Time,
	}

	return message, nil
}

func (r *messageRepo) IsMessageExist(ctx context.Context, msg domain.CreateMessageRequest) (bool, error) {
	var exists bool

	query := `
    SELECT EXISTS (
        SELECT 1
        FROM messages
        WHERE content = $1
    )
    `

	err := r.db.QueryRowContext(ctx, query, msg.Content).Scan(&exists)
	if err != nil {
		logger.Errorf("error checking if task exists: %v", err)
		return false, err
	}

	return exists, nil
}

func (r *messageRepo) GetMessageById(ctx context.Context, id int) (*entity.Message, error) {
	return &entity.Message{}, nil
}

func (r *messageRepo) GetMessagesList(ctx context.Context, inp domain.MessagesListParams) ([]*entity.Message, error) {
	return nil, nil
}

func (r *messageRepo) DeleteMessageById(ctx context.Context, id int) error {
	return nil
}

// func (r *MessageRepo) UpdateMessage(ctx context.Context, msg domain.UpdateMessageRequest) (*entity.Message, error) {

// }
