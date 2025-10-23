ALTER TABLE userWorkout
    ADD COLUMN userId INT NULL REFERENCES users(id);