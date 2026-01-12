package loadprescription

import (
	"database/sql"
	"fmt"

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

func (controller *Controller) CreateLoadPrescriptionForWorkoutSet(
	set *types.TemplateExerciseSet,
	loadPrescriptionStmt *sql.Stmt,
	fixLoadPrescriptionStmt *sql.Stmt,
	percentageOneRepMaxStmt *sql.Stmt,
	rpeStmt *sql.Stmt,
) (int, error) {

	if !controller.IsValidLoadPrescriptionType(set.LoadPresciptionType) {
		return -1, fmt.Errorf("Invalid load type ID found: %v", set.LoadPresciptionType)
	}

	// Create loadPrescription
	loadPrescriptionId, err := controller.CreateLoadPrescriptionFromStatement(
		loadPrescriptionStmt,
		*set.LoadPresciptionType,
	)
	if err != nil {
		return -1, err
	}

	// Insert to the correct table
	switch *set.LoadPresciptionType {
	case types.FIXED:
		err = controller.CreateFixedLoadPrescriptionFromStatement(fixLoadPrescriptionStmt, loadPrescriptionId, set.FixedLoadPrescription)
	case types.PERCENTAGE_ONE_REP_MAX:
		err = controller.CreatePercentageMaxLoadPrescriptionFromStatement(percentageOneRepMaxStmt, loadPrescriptionId, set.PercentageOneRepMaxLoadPrescription)
	case types.RPE:
		err = controller.CreateRPELoadPrescriptionFromStatement(rpeStmt, loadPrescriptionId, set.RPELoadPrescription)
	default:
		err = fmt.Errorf("Non implemented type found. Found %v", *set.LoadPresciptionType)
	}
	if err != nil {
		return -1, err
	}

	return loadPrescriptionId, nil
}

func (controller *Controller) CreateLoadPrescriptionFromStatement(
	stmt *sql.Stmt,
	loadPrescriptionType types.LoadPresciptionType,
) (int, error) {
	var presscriptionId int

	err := stmt.
		QueryRow(loadPrescriptionType).
		Scan(&presscriptionId)
	if err != nil {
		return -1, err
	}

	return presscriptionId, nil
}

func (controller *Controller) CreateFixedLoadPrescriptionFromStatement(
	stmt *sql.Stmt,
	loadPrescriptionId int,
	fixedLoadPrescription *types.FixedLoadPrescription,
) error {

	if fixedLoadPrescription == nil {
		return fmt.Errorf("Fixed load prescription should not be nil")
	}

	_, err := stmt.
		Exec(
			loadPrescriptionId,
			fixedLoadPrescription.Weight,
		)
	return err
}

func (controller *Controller) CreatePercentageMaxLoadPrescriptionFromStatement(
	stmt *sql.Stmt,
	loadPrescriptionId int,
	percentageLoadPrescription *types.PercentageOneRepMaxLoadPrescription,
) error {

	if percentageLoadPrescription == nil {
		return fmt.Errorf("Percentage load prescription should not be nil")
	}

	_, err := stmt.
		Exec(
			loadPrescriptionId,
			percentageLoadPrescription.Percentage,
		)
	return err
}

func (controller *Controller) CreateRPELoadPrescriptionFromStatement(
	stmt *sql.Stmt,
	loadPrescriptionId int,
	rpeLoadPrescription *types.RPELoadPrescription,
) error {

	if rpeLoadPrescription == nil {
		return fmt.Errorf("rpe load prescription should not be nil")
	}

	_, err := stmt.
		Exec(
			loadPrescriptionId,
			rpeLoadPrescription.RPE,
		)
	return err
}

func (controller *Controller) CreateLoadPrescriptionStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare("INSERT INTO load_prescription (type_id) VALUES ($1) RETURNING id")
}

func (controller *Controller) CreateFixedLoadPrescriptionStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare("INSERT INTO fixed_load_prescription (id, weight) VALUES ($1, $2)")
}

func (controller *Controller) CreatePercentageOneRepMaxLoadPrescriptionStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare("INSERT INTO percentage_one_rep_max_load_prescription (id, percentage) VALUES ($1, $2)")
}

func (controller *Controller) CreateRPELoadPrescriptionStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare("INSERT INTO rpe_load_prescription (id, rpe) VALUES ($1, $2)")
}
