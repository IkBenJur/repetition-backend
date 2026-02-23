ALTER TABLE userworkoutexercise RENAME TO user_workout_exercise;

ALTER TABLE user_workout_exercise ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE user_workout_exercise RENAME COLUMN userworkoutid TO user_workout_id;
ALTER TABLE user_workout_exercise RENAME COLUMN exerciseid TO exercise_id;
ALTER TABLE user_workout_exercise RENAME COLUMN createdat TO created_at;
