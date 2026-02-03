package gomigrations

import (
	"context"

	"github.com/IkBenJur/repetition-backend/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const FixedLoadPrescriptionType = 0

func MigrateVersion20LoadPrescription() error {
	db, err := config.ConnectDatabase(config.Envs)
	if err != nil {
		return err
	}

	sets, err := findAllSets(db)
	if err != nil {
		return err
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, set := range sets {
		var addedPrescriptionId int

		err := db.QueryRow(
			context.Background(),
			"INSERT INTO load_prescription (type_id) VALUES ($1) RETURNING id",
			FixedLoadPrescriptionType,
		).Scan(&addedPrescriptionId)
		if err != nil {
			return err
		}

		_, err = db.Exec(
			context.Background(),
			"INSERT INTO fixed_load_prescription (id, weight) VALUES ($1, $2)",
			addedPrescriptionId,
			set.weight,
		)
		if err != nil {
			return err
		}

		_, err = db.Exec(
			context.Background(),
			"UPDATE userworkoutexerciseset SET load_prescription_id = $1 WHERE id = $2",
			addedPrescriptionId,
			set.id,
		)
		if err != nil {
			return err
		}

	}

	return tx.Commit(context.Background())
}

// Only fields relevant for migration
type set struct {
	id     int
	weight float64
}

func findAllSets(db *pgxpool.Pool) ([]set, error) {
	sets := make([]set, 0)

	rows, err := db.Query(context.Background(), "SELECT id, weight FROM userworkoutexerciseset")
	if err != nil {
		return sets, err
	}
	defer rows.Close()

	for rows.Next() {
		var set set

		err := rows.Scan(&set.id, &set.weight)
		if err != nil {
			return sets, err
		}

		sets = append(sets, set)
	}

	return sets, nil
}
