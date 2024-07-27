package repository

import "github.com/jmoiron/sqlx"

type Repository interface{}

type repoPostgres struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *repoPostgres {
	return &repoPostgres{db: db}
}
