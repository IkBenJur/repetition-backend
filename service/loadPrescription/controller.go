package loadprescription

import (
	"database/sql"

	"github.com/IkBenJur/repetition-backend/types"
)

type Controller struct {
	db *sql.DB
}

func NewController(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

func (controller *Controller) IsValidLoadPrescriptionType(loadTypeId *types.LoadPresciptionType) bool {
	if loadTypeId == nil {
		return false
	}

	if *loadTypeId == types.FIXED {
		return true
	}

	if *loadTypeId == types.PERCENTAGE_ONE_REP_MAX {
		return true
	}

	if *loadTypeId == types.RPE {
		return true
	}

	return false
}

func (controller *Controller) CreateLoadPrescriptionStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare("INSERT INTO load_prescription (type_id) VALUES ($1) RETURNING id")
}

func (controller *Controller) CreateFixedLoadPrescriptionStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare("INSERT INTO fixed_load_prescription (id, weight) VALUES ($1, $2)")
}
