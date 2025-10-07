package main

import (
	"errors"
	"log"
	"os"

	"github.com/IkBenJur/repetition-backend/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	envs := config.InitConfig()

	m, err := migrate.New("file://"+"./migrate/migrations",envs.DbConnectionUrl)
	if err != nil {
		log.Fatalf("Error connectiong migration: %v", err)
		return
	}
	
	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Error running migrations: %v", err)
			return;
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Error running migrations: %v", err)
			return;
		}
	}
	

	srcErr, dbErr := m.Close()
	if srcErr != nil {
		log.Fatalf("Src error: %v", srcErr)
	}

	if dbErr != nil {
		log.Fatalf("Db error: %v", dbErr)
	}

	log.Print("Migrations applied")
}