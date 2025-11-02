package api

import "main/database"

type RateLimitConfig struct {
	MaxRequestPerMin int64
}

type Config struct {
	DatabaseConfig  database.DatabaseConfig
	RateLimitConfig RateLimitConfig
}