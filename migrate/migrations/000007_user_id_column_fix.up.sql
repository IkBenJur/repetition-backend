UPDATE userworkout
    SET userid = 1;

ALTER TABLE userWorkout 
    ALTER COLUMN userId SET NOT NULL;