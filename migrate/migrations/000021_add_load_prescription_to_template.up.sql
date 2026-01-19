-- No need for backfill. Templates not in production yet
ALTER TABLE template_exercise_set
    ADD COLUMN load_prescription_id INT REFERENCES load_prescription(id);

ALTER TABLE template_exercise_set
    DROP COLUMN weight_goal;
