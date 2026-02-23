-- Forgot to add in previous migrations
ALTER TABLE user_workout_exercise ADD COLUMN user_id INT NOT NULL DEFAULT 1 REFERENCES user_table(id);

ALTER TABLE userworkoutexerciseset RENAME TO user_workout_exercise_set;

ALTER TABLE user_workout_exercise_set ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE user_workout_exercise_set ADD COLUMN user_id INT NOT NULL DEFAULT 1 REFERENCES user_table(id);
ALTER TABLE user_workout_exercise_set RENAME COLUMN createdat TO created_at;
ALTER TABLE user_workout_exercise_set RENAME COLUMN userworkoutexerciseid TO user_workout_exercise_id;
