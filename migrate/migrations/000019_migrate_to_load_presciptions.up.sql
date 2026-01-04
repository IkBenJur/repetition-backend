-- This migration consists of a Go migration file as well
-- New load prescription of the FIXED type will be created using existing set data
ALTER TABLE userworkoutexerciseset
    ADD COLUMN load_prescription_id INT REFERENCES load_prescription(id);
