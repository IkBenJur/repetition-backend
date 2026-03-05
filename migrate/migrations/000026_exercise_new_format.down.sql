ALTER TABLE exercise DROP COLUMN updated_at;
ALTER TABLE exercise RENAME COLUMN created_at TO createdAt;
ALTER TABLE exercise RENAME COLUMN muscle_group TO muscleGroup;

