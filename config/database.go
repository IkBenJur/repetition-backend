package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
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

func ConnectDatabase(config Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.DbConnectionUrl)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	return db, nil
}
