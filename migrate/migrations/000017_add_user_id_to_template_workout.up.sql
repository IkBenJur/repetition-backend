ALTER TABLE workout_template ADD COLUMN user_id INT NOT NULL REFERENCES users(id);
