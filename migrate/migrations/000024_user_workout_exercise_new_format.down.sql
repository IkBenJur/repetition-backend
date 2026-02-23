ALTER TABLE user_workout_exercise RENAME TO userworkoutexercise;

ALTER TABLE userworkoutexercise DROP COLUMN updated_at;
ALTER TABLE userworkoutexercise RENAME COLUMN user_workout_id TO userworkoutid;
ALTER TABLE userworkoutexercise RENAME COLUMN exercise_id TO exerciseid;
ALTER TABLE userworkoutexercise RENAME COLUMN created_at TO createdat;
