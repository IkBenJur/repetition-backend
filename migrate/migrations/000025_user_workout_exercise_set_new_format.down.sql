ALTER TABLE user_workout_exercise DROP COLUMN user_id;

ALTER TABLE user_workout_exercise_set RENAME TO userworkoutexerciseset;

ALTER TABLE userworkoutexerciseset DROP COLUMN updated_at;
ALTER TABLE userworkoutexerciseset DROP COLUMN user_id;
ALTER TABLE userworkoutexerciseset RENAME COLUMN created_at TO createdat;
ALTER TABLE userworkoutexerciseset RENAME COLUMN user_workout_exercise_id TO userworkoutexerciseid;
