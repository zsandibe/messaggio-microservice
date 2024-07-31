package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zsandibe/messaggio-microservice/internal/entity"
	logger "github.com/zsandibe/messaggio-microservice/pkg"
)

type statisticRepo struct {
	db *sqlx.DB
}

func NewStatisticRepo(db *sqlx.DB) *statisticRepo {
	return &statisticRepo{db: db}
}

func (r *statisticRepo) GetStatsList(ctx context.Context) ([]*entity.Stats, error) {
	var stats []*entity.Stats
	query := `
		SELECT id,processed_count,
		last_processed_message,
		updated_at
		FROM statistics
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error("error querying statistics: %v", err)
		return nil, err
	}

	for rows.Next() {
		var stat *entity.Stats
		err = rows.Scan(&stat.Id, &stat.ProcessedCount,
			&stat.LastProcessedMessageId, &stat.UpdatedAt)

		if err != nil {
			logger.Error("error scaning row: %v", err)
			return nil, err
		}

		stats = append(stats, stat)
	}
	return stats, nil
}

func (r *statisticRepo) GetStatById(ctx context.Context, id int) (*entity.Stats, error) {
	var stat entity.Stats
	var updatedAt sql.NullTime

	query := `
		SELECT m.id,m.content,
		m.status,m.created_at,
		m.processed_at
		FROM messages m 
		WHERE m.id = $1
	`

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stat.Id,
		&stat.ProcessedCount,
		&stat.LastProcessedContent,
		&stat.LastProcessedMessageId,
		&updatedAt,
	); err != nil {
		return &stat, err
	}

	stat.UpdatedAt = updatedAt.Time

	return &stat, nil
}
