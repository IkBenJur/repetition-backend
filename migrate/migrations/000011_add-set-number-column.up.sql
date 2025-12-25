ALTER TABLE userWorkoutExerciseSet ADD COLUMN IF NOT EXISTS set_number INT;

-- Update the table using a table view of each set ID with its own set number. Partitioned by each exerciseId
UPDATE userWorkoutExerciseSet uwes
SET set_number = set_number_view.set_number
FROM (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY userworkoutexerciseid
            ORDER BY createdat, id
        ) AS set_number
    FROM userWorkoutExerciseSet
) set_number_view
WHERE uwes.id = set_number_view.id;

ALTER TABLE
    userWorkoutExerciseSet
ALTER COLUMN
    set_number SET NOT NULL;

ALTER TABLE
    userWorkoutExerciseSet
ADD CONSTRAINT
    unique_set_number_exercise_id
UNIQUE
    (userworkoutexerciseid, set_number);
