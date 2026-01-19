package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/IkBenJur/repetition-backend/config"
	gomigrations "github.com/IkBenJur/repetition-backend/migrate/go_migrations"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Version numbers Go migrations files
const (
	loadPrescriptionVersionNumber = 20
	RunAll                        = -1
)

func main() {
	envs := config.InitConfig()

	m, err := migrate.New("file://"+"./migrate/migrations", envs.DbConnectionUrl)
	if err != nil {
		log.Fatalf("Error connecting migration: %v", err)
	}
	defer closeMigrate(m)

	// Args always contains program name so at minimum 1.
	if len(os.Args) == 1 {
		log.Fatal("Should contain at least one argument: up, down or migrate <version_number>")
	}

	// Depending on the version number we might need to run a Go migration file
	currentVersionNumber, err := getCurrentVersion(m)
	if err != nil {
		log.Fatalf("Could't get current version number: %v", err)
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		handleUp(m, currentVersionNumber)
	case "down":
		handleDown(m)
	case "migrate":
		if len(os.Args) < 3 {
			log.Fatal("migrate command requires a version number")
		}
		targetVersion, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		handleMigrate(m, targetVersion, currentVersionNumber)
	default:
		log.Fatalf("Unknown command: %s. Use 'up', 'down', or 'migrate <version>'", cmd)
	}

	log.Println("Migrations applied")
}

func handleUp(m *migrate.Migrate, currentVersion uint) {
	// Run the Go migration file inbetween versions when currentVersionNumber is less
	if currentVersion < uint(loadPrescriptionVersionNumber) {

		// Run to the specific version when the Go migration should run
		if err := runMigration(loadPrescriptionVersionNumber, m); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Error running migrations: %v", err)
		}

		if err := gomigrations.MigrateVersion20LoadPrescription(); err != nil {
			log.Fatalf("Error running GO migration V19 file: %v", err)
		}

	}

	// Run the rest of the migrations
	if err := runMigration(RunAll, m); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Error running migrations: %v", err)
	}
}

func handleDown(m *migrate.Migrate) {
	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Error running migrations: %v", err)
	}
}

func handleMigrate(m *migrate.Migrate, inputVersion int, currentVersion uint) {
	migrationUp := inputVersion > int(currentVersion)
	beforeGoMigration := currentVersion < loadPrescriptionVersionNumber
	crossesGoMigration := inputVersion >= loadPrescriptionVersionNumber

	// Run the Go migration file if going up and crossing version 19
	if migrationUp && beforeGoMigration && crossesGoMigration {

		// Run to the specific version when the Go migration should run
		if err := runMigration(loadPrescriptionVersionNumber, m); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Error running migrations: %v", err)
		}

		if err := gomigrations.MigrateVersion20LoadPrescription(); err != nil {
			log.Fatalf("Error running GO migration V19 file: %v", err)
		}

	}

	// Run the rest of the migrations
	if err := runMigration(inputVersion, m); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Error running migrations: %v", err)
	}
}

func runMigration(targetVersion int, m *migrate.Migrate) error {
	// When target version is not set run the up command
	if targetVersion == -1 {
		return m.Up()
	}

	// Otherwise run to specific version
	return m.Migrate(uint(targetVersion))
}

func getCurrentVersion(m *migrate.Migrate) (uint, error) {
	currentVersionNumber, _, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, err
	}

	// When version number isn't set yet
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, nil
	}

	return currentVersionNumber, nil
}

func closeMigrate(m *migrate.Migrate) {
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		log.Printf("Source close error: %v", srcErr)
	}
	if dbErr != nil {
		log.Printf("Database close error: %v", dbErr)
	}
}
