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
	var status bool = false
	var processedAt sql.NullTime
	startTime := time.Now()

	query := `
	INSERT INTO messages (content,status,created_at,processed_at) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, msg.Content, status, startTime, processedAt).Scan(&id)
	if err != nil {
		logger.Errorf("Error in inserting message: %v", err)
		return &entity.Message{}, domain.ErrCreatingMessage
	}

	message = &entity.Message{
		Id:          id,
		Content:     msg.Content,
		IsProcessed: status,
		CreatedAt:   startTime,
		ProcessedAt: processedAt.Time,
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
	var processedAt sql.NullTime

	query := `
		SELECT m.id,m.content,
		m.status,m.created_at,
		m.processed_at
		FROM messages m 
		WHERE m.id = $1
	`

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.Id,
		&message.Content,
		&message.IsProcessed,
		&message.CreatedAt,
		&processedAt,
	); err != nil {
		return &message, err
	}

	message.ProcessedAt = processedAt.Time

	return &message, nil
}

func (r *messageRepo) GetMessagesList(ctx context.Context, params domain.MessagesListParams) ([]*entity.Message, error) {
	messages := make([]*entity.Message, 0)
	fmt.Println(params)
	var (
		args    []interface{}
		where   []string
		orderBy []string
	)

	if params.Content != "" {
		where = append(where, "content = $"+strconv.Itoa(len(args)+1))
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
		var processedAt sql.NullTime
		msg := &entity.Message{}
		if err := rows.Scan(&msg.Id, &msg.Content, &msg.IsProcessed, &msg.CreatedAt, &processedAt); err != nil {
			logger.Error(err)
			return nil, err
		}

		msg.ProcessedAt = processedAt.Time

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

func (r *messageRepo) UpdateStatus(ctx context.Context, id int, content string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("error beginning transaction: %v", err))
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	queryUpdateMessage := `
        UPDATE messages 
        SET status = $1, processed_at = NOW()
        WHERE id = $2
    `
	flag := true
	_, err = tx.ExecContext(ctx, queryUpdateMessage, flag, id)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("error executing query: %v", err))
		return err
	}

	var statsID int
	var processedCount int

	querySelectStats := `
        SELECT id, processed_count
        FROM message_stats
        WHERE last_processed_message_content = $1
    `

	row := tx.QueryRowContext(ctx, querySelectStats, content)
	err = row.Scan(&statsID, &processedCount)
	fmt.Println(content, id)
	if err != nil {
		if err == sql.ErrNoRows {

			queryInsertStats := `
                INSERT INTO message_stats (processed_count,last_processed_message_content, last_processed_message_id, updated_at)
                VALUES (1,$1, $2, NOW())
            `
			_, err = tx.ExecContext(ctx, queryInsertStats, content, id)
			if err != nil {
				logger.Error(ctx, fmt.Errorf("error inserting into message_stats: %v", err))
				return err
			}
		} else {
			logger.Error(ctx, fmt.Errorf("error selecting from message_stats: %v", err))
			return err
		}
	} else {

		queryUpdateStats := `
            UPDATE message_stats
            SET processed_count = $1, updated_at = NOW()
            WHERE id = $2
        `
		_, err = tx.ExecContext(ctx, queryUpdateStats, processedCount+1, statsID)
		if err != nil {
			logger.Error(ctx, fmt.Errorf("error updating message_stats: %v", err))
			return err
		}
	}

	return nil
}
