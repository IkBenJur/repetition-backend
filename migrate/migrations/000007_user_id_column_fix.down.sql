UPDATE userworkout
    SET userid = null;

ALTER TABLE userWorkout 
        ALTER COLUMN userId DROP NOT NULL;