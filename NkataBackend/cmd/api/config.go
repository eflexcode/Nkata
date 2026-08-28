package api

import "main/database"

type RateLimitConfig struct {
	MaxRequestPerMin int64
}

type RedisConfig struct{

	Addre string
	Password string
	Db int

}

type Config struct {
	DatabaseConfig  database.DatabaseConfig
	RateLimitConfig RateLimitConfig
	RedisConfig RedisConfig
}