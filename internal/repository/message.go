package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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
	var message entity.Message

	query := `
		SELECT m.id,m.content,
		m.status,m.created_at,
		m.updated_at
		FROM messages m 
		WHERE m.id = $1
	`

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.Id,
		&message.Content,
		&message.IsProcessed,
		&message.CreatedAt,
		&message.ProcessedAt,
	); err != nil {
		return &message, err
	}
	return &message, nil
}

func (r *messageRepo) GetMessagesList(ctx context.Context, params domain.MessagesListParams) ([]*entity.Message, error) {
	messages := make([]*entity.Message, 0)
	var (
		args    []interface{}
		where   []string
		orderBy []string
	)

	if params.Content != "" {
		where = append(where, "passport_serie = $"+strconv.Itoa(len(args)+1))
		args = append(args, params.Content)
	}

	query := "SELECT * FROM messages"
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	if len(orderBy) > 0 {
		query += " ORDER BY " + strings.Join(orderBy, ", ")
	}
	if params.Limit > 0 {
		query += " LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, params.Limit)
	}
	if params.Offset > 0 {
		query += " OFFSET $" + strconv.Itoa(len(args)+1)
		args = append(args, params.Offset)
	}

	fmt.Println(query)
	fmt.Println(args)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg *entity.Message
		if err := rows.Scan(&msg.Id, &msg.Content, &msg.IsProcessed, &msg.CreatedAt, &msg.ProcessedAt); err != nil {
			logger.Error(err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		logger.Error(err)
		return nil, err
	}

	return messages, nil
}

func (r *messageRepo) DeleteMessageById(ctx context.Context, id int) error {
	query := `
		DELETE FROM messages WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("error with executing query: %v", err))
		return err
	}

	return nil
}

func (r *messageRepo) UpdateStatus(ctx context.Context, id int) error {
	query := `
		UPDATE messages 
		SET status = $1
		WHERE id = $2
	`

	flag := true

	_, err := r.db.ExecContext(ctx, query, flag, id)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("error with executing query: %v", err))
		return err
	}
	return nil
}
