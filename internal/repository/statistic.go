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
		last_processed_message_content,
		last_processed_message_id,
		updated_at
		FROM message_stats
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error("error querying statistics: %v", err)
		return nil, err
	}

	for rows.Next() {
		var stat entity.Stats
		err = rows.Scan(&stat.Id, &stat.ProcessedCount,
			&stat.LastProcessedContent,
			&stat.LastProcessedMessageId, &stat.UpdatedAt)

		if err != nil {
			logger.Error("error scaning row: %v", err)
			return nil, err
		}

		stats = append(stats, &stat)
	}
	return stats, nil
}

func (r *statisticRepo) GetStatById(ctx context.Context, id int) (*entity.Stats, error) {
	var stat entity.Stats
	var updatedAt sql.NullTime

	query := `
		SELECT s.id,s.processed_count,
		s.last_processed_message_content,s.last_processed_message_id,
		s.updated_at
		FROM message_stats s 
		WHERE s.id = $1
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
