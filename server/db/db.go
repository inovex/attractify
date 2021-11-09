package db

import (
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *DB {
	return &DB{db: db}
}

type DB struct {
	db *sqlx.DB
}

func (db *DB) Close() error {
	return db.db.Close()
}
