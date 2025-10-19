CREATE TABLE IF NOT EXISTS exercise (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    muscleGroup VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO exercise (name, muscleGroup) VALUES
    ('Squat', 'Legs'),
    ('Bench press', 'Chest'),
    ('Deadlift', 'Legs'),
    ('Barbell row', 'Back');