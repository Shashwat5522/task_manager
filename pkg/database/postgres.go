package database

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	// TODO: Implement database connection
	return nil, nil
}
