package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DbConnectionUrl string
	JWTSecret string
}

func InitConfig() Config {
	dsn := os.Getenv("DB_CONN_URL")
	JWTSecret := os.Getenv("JWT_SECRET")
	return Config {
		DbConnectionUrl: dsn,
		JWTSecret: JWTSecret,
	}
}

var Envs = InitConfig()

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