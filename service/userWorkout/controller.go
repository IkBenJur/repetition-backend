package userWorkout

import (
	"database/sql"
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

	exerciseStmt, err := tx.Prepare("INSERT INTO userworkoutexercise (userworkoutid, exerciseid) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return -1, err
	}
	defer exerciseStmt.Close()

	exerciseSetStmt, err := tx.Prepare("INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight, set_number) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return -1, err
	}
	defer exerciseSetStmt.Close()

	for _, exercise := range workout.UserWorkoutExercises {
		var exerciseId int
		err = exerciseStmt.QueryRow(workoutId, exercise.ExerciseId).Scan(&exerciseId)
		if err != nil {
			return -1, err
		}

		for _, set := range exercise.UserWorkoutExerciseSets {
			_, err := exerciseSetStmt.Exec(exerciseId, set.Reps, set.Weight, set.SetNumber)
			if err != nil {
				return -1, err
			}
		}
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

	userId := 0
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			return 0, err
		}
	}

	return userId, nil
}

func (controller *Controller) FindActiveWorkoutForUserId(id int) (*types.UserWorkout, error) {
	var userWorkout *types.UserWorkout
	exerciseMap := map[int]*types.UserWorkoutExercise{}

	rows, err := controller.db.Query(`SELECT
										uw.id, uw.name, uw.datestart, uw.dateend, uw.createdat, uw.userid,
										uwe.id, uwe.userworkoutid, uwe.exerciseid, exer.name, uwe.createdat,
										uwes.id, uwes.userworkoutexerciseid, uwes.reps, uwes.weight, uwes.set_number, uwes.createdat
									FROM userworkout uw
									LEFT JOIN userworkoutexercise uwe
										ON uw.id = uwe.userworkoutid
									LEFT JOIN userworkoutexerciseset uwes
										ON uwe.id = uwes.userworkoutexerciseid
									LEFT JOIN exercise exer
										ON exer.id = uwe.exerciseid
									WHERE uw.id = (
										SELECT active_userworkout_id
										FROM users
										WHERE id = $1
									) ORDER BY uwes.set_number
	`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			uw   types.UserWorkout
			uwe  types.UserWorkoutExercise
			uwes types.UserWorkoutExerciseSet
		)

		err := rows.Scan(
			&uw.ID, &uw.Name, &uw.DateStart, &uw.DateEnd, &uw.CreatedAt, &uw.UserId,
			&uwe.ID, &uwe.UserWorkoutId, &uwe.ExerciseId, &uwe.ExerciseName, &uwe.CreatedAt,
			&uwes.ID, &uwes.UserWorkoutExerciseId, &uwes.Reps, &uwes.Weight, &uwes.SetNumber, &uwes.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if userWorkout == nil {

			//Initialize an empty array
			uw.UserWorkoutExercises = make([]*types.UserWorkoutExercise, 0)

			userWorkout = &uw
		}

		uweIsNotNull := uwe.ID != 0
		if uweIsNotNull {
			// Check whether exercise had already been added or not
			// If not add it to the map and workout
			if _, ok := exerciseMap[uwe.ID]; !ok {

				// Initialize an empty array
				uwe.UserWorkoutExerciseSets = make([]*types.UserWorkoutExerciseSet, 0)

				exerciseMap[uwe.ID] = &uwe
				userWorkout.UserWorkoutExercises = append(userWorkout.UserWorkoutExercises, &uwe)

			}
		}

		uwesIsNotNull := uwes.ID != 0
		if uwesIsNotNull {
			// Add the set to the existing exercise
			parent := exerciseMap[uwe.ID]
			parent.UserWorkoutExerciseSets = append(parent.UserWorkoutExerciseSets, &uwes)
		}
	}

	return userWorkout, nil
}

func (controller *Controller) MarkUserWorkoutAsComplete(id int) error {
	_, err := controller.db.Exec("UPDATE userworkout SET dateend = $1 WHERE id = $2", time.Now(), id)
	return err
}
