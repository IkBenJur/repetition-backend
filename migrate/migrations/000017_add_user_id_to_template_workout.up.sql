ALTER TABLE workout_template ADD COLUMN user_id INT NOT NULL DEFAULT 1 REFERENCES users(id);

ALTER TABLE workout_template ALTER COLUMN user_id DROP DEFAULT;
