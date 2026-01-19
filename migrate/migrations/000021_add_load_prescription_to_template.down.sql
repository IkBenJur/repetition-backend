ALTER TABLE template_exercise_set
    DROP COLUMN load_prescription_id;

ALTER TABLE template_exercise_set
    ADD COLUMN weight_goal DECIMAL(6,3);
