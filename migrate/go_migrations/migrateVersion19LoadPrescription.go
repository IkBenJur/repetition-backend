package gomigrations

import (
	"database/sql"

	"github.com/IkBenJur/repetition-backend/config"
)

const FixedLoadPrescriptionType = 0

func MigrateVersion19LoadPrescription() error {
	db, err := config.ConnectDatabase(config.Envs)
	if err != nil {
		return err
	}

	sets, err := findAllSets(db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertLoadPrescriptionStatement, err := tx.Prepare("INSERT INTO load_prescription (type_id) VALUES ($1) RETURNING id")
	if err != nil {
		return err
	}
	defer insertLoadPrescriptionStatement.Close()

	insertFixedLoadPrescriptionStatement, err := tx.Prepare("INSERT INTO fixed_load_prescription (id, weight) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer insertFixedLoadPrescriptionStatement.Close()

	updateSetStatement, err := tx.Prepare("UPDATE userworkoutexerciseset SET load_prescription_id = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	defer updateSetStatement.Close()

	for _, set := range sets {
		var addedPrescriptionId int

		err := insertLoadPrescriptionStatement.QueryRow(FixedLoadPrescriptionType).Scan(&addedPrescriptionId)
		if err != nil {
			return err
		}

		_, err = insertFixedLoadPrescriptionStatement.Exec(addedPrescriptionId, set.weight)
		if err != nil {
			return err
		}

		_, err = updateSetStatement.Exec(addedPrescriptionId, set.id)
		if err != nil {
			return err
		}

	}

	return tx.Commit()
}

// Only fields relevant for migration
type set struct {
	id     int
	weight float64
}

func findAllSets(db *sql.DB) ([]set, error) {
	sets := make([]set, 0)

	rows, err := db.Query("SELECT id, weight FROM userworkoutexerciseset")
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
