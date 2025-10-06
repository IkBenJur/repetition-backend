package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDatabase(){
	var err error
	dsn := os.Getenv("DB_CONN_URL")
    DB, err = sql.Open("pgx", dsn)
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
	
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}
    fmt.Println("Database connected!")
}