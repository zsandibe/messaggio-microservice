package repository

import "github.com/jmoiron/sqlx"

type statisticRepo struct {
	db *sqlx.DB
}

func NewStatisticRepo(db *sqlx.DB) *statisticRepo {
	return &statisticRepo{db: db}
}
