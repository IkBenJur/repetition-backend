ALTER TABLE userWorkoutExercise ADD COLUMN IF NOT EXISTS exercise_number INT;

-- Update the table using a table view of each exercise ID with its own exercise number. Partitioned by each workoutOutId
UPDATE userWorkoutExercise uwe
SET exercise_number = exercise_number_view.exercise_number
FROM (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY userworkoutid
            ORDER BY createdat, id
        ) AS exercise_number
    FROM userWorkoutExercise
) exercise_number_view
WHERE uwe.id = exercise_number_view.id;

ALTER TABLE
    userWorkoutExercise
ALTER COLUMN
    exercise_number SET NOT NULL;

ALTER TABLE
    userWorkoutExercise
ADD CONSTRAINT
    unique_exercise_number_workout_id
UNIQUE
    (userWorkoutId, exercise_number);
