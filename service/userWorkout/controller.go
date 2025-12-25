package userWorkout

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/IkBenJur/repetition-backend/types"
)

type Controller struct {
	db *sql.DB
}

func NewController(db *sql.DB) *Controller {
	return &Controller{db: db}
}

func (controller *Controller) CreateNewUserWorkout(workout types.UserWorkout) (int, error) {
	tx, err := controller.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	var workoutId int
	err = tx.QueryRow("INSERT INTO userWorkout (name, userId) VALUES ($1, $2) RETURNING id", workout.Name, workout.UserId).Scan(&workoutId)
	if err != nil {
		return -1, err
	}

	exerciseStmt, err := tx.Prepare("INSERT INTO userworkoutexercise (userworkoutid, exerciseid, exercise_number) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return -1, err
	}
	defer exerciseStmt.Close()

	exerciseSetStmt, err := tx.Prepare("INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight, set_number) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return -1, err
	}
	defer exerciseSetStmt.Close()

	exerciseNumber := 1
	for _, exercise := range workout.UserWorkoutExercises {
		var exerciseId int

		// Create new varaible for specific exercise
		newExerciseNumber := exerciseNumber
		exercise.ExerciseNumber = &newExerciseNumber

		err = exerciseStmt.QueryRow(workoutId, exercise.ExerciseId, exercise.ExerciseNumber).Scan(&exerciseId)
		if err != nil {
			return -1, err
		}

		setNumber := 1
		for _, set := range exercise.UserWorkoutExerciseSets {

			// Create new varaible for specific set
			newSetNumber := setNumber
			set.SetNumber = &newSetNumber

			_, err := exerciseSetStmt.Exec(exerciseId, set.Reps, set.Weight, set.SetNumber)
			if err != nil {
				return -1, err
			}

			setNumber++
		}

		exerciseNumber++
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return workoutId, nil
}

func (controller *Controller) FindUserIdForUserworkoutId(id int) (int, error) {
	rows, err := controller.db.Query("SELECT userid FROM userworkout WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	userId := 0
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return 0, err
		}
	}

	return userId, nil
}

func (controller *Controller) FindAllWorkoutsForUserId(userId int) ([]*types.UserWorkout, error) {
	//Initialize an empty array
	userWorkouts := make([]*types.UserWorkout, 0)
	rows, err := controller.db.Query(`
		SELECT id, userId, name, dateStart, dateEnd, createdAt
		FROM userworkout
		WHERE userId = $1 ORDER BY id DESC`, userId)
	if err != nil {
		return userWorkouts, err
	}
	defer rows.Close()

	for rows.Next() {
		var userWorkout types.UserWorkout

		err := rows.Scan(&userWorkout.ID, &userWorkout.UserId, &userWorkout.Name, &userWorkout.DateStart, &userWorkout.DateEnd, &userWorkout.CreatedAt)
		if err != nil {
			return userWorkouts, err
		}

		userWorkouts = append(userWorkouts, &userWorkout)
	}

	return userWorkouts, nil
}

func (controller *Controller) FindActiveWorkoutForUserId(id int) (*types.UserWorkout, error) {
	whereClause := `uw.id = (
		SELECT active_userworkout_id
		FROM users
		WHERE id = $1
	)`
	return controller.findWorkout(whereClause, id)
}

func (controller *Controller) FindById(workoutId int) (*types.UserWorkout, error) {
	whereClause := `uw.id = $1`
	return controller.findWorkout(whereClause, workoutId)
}

func (controller *Controller) findWorkout(whereClause string, id int) (*types.UserWorkout, error) {
	var userWorkout *types.UserWorkout
	exerciseMap := map[int]*types.UserWorkoutExercise{}

	query := fmt.Sprintf(`SELECT
					uw.id, uw.name, uw.datestart, uw.dateend, uw.createdat, uw.userid,
					uwe.id, uwe.userworkoutid, uwe.exerciseid, exer.name, uwe.createdat,
					uwes.id, uwes.userworkoutexerciseid, uwes.reps, uwes.weight, uwes.set_number, uwes.is_done, uwes.createdat
				FROM userworkout uw
				LEFT JOIN userworkoutexercise uwe
					ON uw.id = uwe.userworkoutid
				LEFT JOIN userworkoutexerciseset uwes
					ON uwe.id = uwes.userworkoutexerciseid
				LEFT JOIN exercise exer
					ON exer.id = uwe.exerciseid
				WHERE %s ORDER BY uwe.exercise_number, uwes.set_number
		`, whereClause)

	rows, err := controller.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// Using left joins so some fields might be empty
		// Use sql.Nullable types
		var (
			uw types.UserWorkout

			// Nullable fields for userWorkoutExercise
			uweID            sql.NullInt64
			uweUserWorkoutId sql.NullInt64
			uweExerciseId    sql.NullInt64
			uweExerciseName  sql.NullString
			uweCreatedAt     sql.NullTime

			// Nullable fields for userWorkoutExerciseSet
			uwesID                    sql.NullInt64
			uwesUserWorkoutExerciseId sql.NullInt64
			uwesReps                  sql.NullInt64
			uwesWeight                sql.NullFloat64
			uwesSetNumber             sql.NullInt64
			uwesIsDone                sql.NullBool
			uwesCreatedAt             sql.NullTime
		)

		err := rows.Scan(
			&uw.ID, &uw.Name, &uw.DateStart, &uw.DateEnd, &uw.CreatedAt, &uw.UserId,
			&uweID, &uweUserWorkoutId, &uweExerciseId, &uweExerciseName, &uweCreatedAt,
			&uwesID, &uwesUserWorkoutExerciseId, &uwesReps, &uwesWeight, &uwesSetNumber, &uwesIsDone, &uwesCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if userWorkout == nil {

			//Initialize an empty array
			uw.UserWorkoutExercises = make([]*types.UserWorkoutExercise, 0)

			userWorkout = &uw
		}

		if uweID.Valid {
			uweIDInt := int(uweID.Int64)
			// Check whether exercise had already been added or not
			// If not add it to the map and workout
			if _, ok := exerciseMap[uweIDInt]; !ok {

				// Construct new exercise
				uwe := types.UserWorkoutExercise{
					ID:                      uweIDInt,
					UserWorkoutId:           int(uweUserWorkoutId.Int64),
					ExerciseId:              int(uweExerciseId.Int64),
					ExerciseName:            &uweExerciseName.String,
					CreatedAt:               uweCreatedAt.Time,
					UserWorkoutExerciseSets: make([]*types.UserWorkoutExerciseSet, 0),
				}

				exerciseMap[uweIDInt] = &uwe
				userWorkout.UserWorkoutExercises = append(userWorkout.UserWorkoutExercises, &uwe)

			}
		}

		if uwesID.Valid {
			uwesIDInt := int(uwesID.Int64)
			uweIDInt := int(uweID.Int64)

			// Contruct new set object
			reps := int(uwesReps.Int64)
			weight := uwesWeight.Float64
			setNumber := int(uwesSetNumber.Int64)
			IsDone := uwesIsDone.Bool

			uwes := types.UserWorkoutExerciseSet{
				ID:                    uwesIDInt,
				UserWorkoutExerciseId: int(uwesUserWorkoutExerciseId.Int64),
				Reps:                  &reps,
				Weight:                &weight,
				SetNumber:             &setNumber,
				IsDone:                IsDone,
				CreatedAt:             uwesCreatedAt.Time,
			}

			// Add the set to the existing exercise
			parent := exerciseMap[uweIDInt]
			parent.UserWorkoutExerciseSets = append(parent.UserWorkoutExerciseSets, &uwes)
		}
	}

	return userWorkout, nil
}

func (controller *Controller) MarkUserWorkoutAsComplete(id int) error {
	_, err := controller.db.Exec("UPDATE userworkout SET dateend = $1 WHERE id = $2", time.Now(), id)
	return err
}
