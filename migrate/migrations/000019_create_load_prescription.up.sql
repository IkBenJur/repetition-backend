CREATE TABLE IF NOT EXISTS load_prescription (
    id SERIAL PRIMARY KEY,
    type_id INT NOT NULL -- Static enum in code
);

-- The different types of load_prescriptions 'share' the load_prescription IDs
CREATE TABLE IF NOT EXISTS fixed_load_prescription (
    id INT PRIMARY KEY REFERENCES load_prescription(id) ON DELETE CASCADE,
    weight DECIMAL(6,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS percentage_one_rep_max_load_prescription (
    id INT PRIMARY KEY REFERENCES load_prescription(id) ON DELETE CASCADE,
    percentage DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rpe_load_prescription (
    id INT PRIMARY KEY REFERENCES load_prescription(id) ON DELETE CASCADE,
    rpe DECIMAL(3,1) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
