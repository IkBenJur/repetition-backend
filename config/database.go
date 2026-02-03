package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
)

type Config struct {
	DbConnectionUrl string
	JWTSecret       string
	FrontEndUrl     string
}

func InitConfig() Config {
	godotenv.Load()

	return Config{
		DbConnectionUrl: getEnv("DATABASE_URL", ""),
		JWTSecret:       getEnv("JWT_SECRET", "a-string-secret-at-least-256-bits-long"),
		FrontEndUrl:     getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

var Envs = InitConfig()

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func ConnectDatabase(config Config) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), config.DbConnectionUrl)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	return dbPool, nil
}
