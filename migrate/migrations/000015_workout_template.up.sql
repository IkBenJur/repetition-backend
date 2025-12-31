CREATE TABLE IF NOT EXISTS workout_template (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS template_workout_exercise (
    id SERIAL PRIMARY KEY,
    workout_template_id INT NOT NULL REFERENCES workout_template(id),
    exercise_id INT NOT NULL REFERENCES exercise(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS template_exercise_set (
    id SERIAL PRIMARY KEY,
    rep_goal INT NOT NULL,
    weight_goal DECIMAL(6,3) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
