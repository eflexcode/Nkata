package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type DatabaseConfig struct {
	Addr         string
	MaxOpenConn  int
	MaxIdealConn int
	MaxIdealTime string
}

type DataRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *DataRepository {
	return &DataRepository{db: db}
}

type PaginatedResponse struct {
	Data       any `json:"data"`
	TotalCount int   `json:"total_count"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
}

func ConnectDatabase(databaseConfig DatabaseConfig) (*sql.DB, error) {

	db, err := sql.Open("postgres", databaseConfig.Addr)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(databaseConfig.MaxOpenConn)
	db.SetMaxIdleConns(databaseConfig.MaxIdealConn)

	duration, err := time.ParseDuration(databaseConfig.MaxIdealTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
