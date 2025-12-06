ALTER TABLE userWorkout 
        ALTER COLUMN userId DROP NOT NULL;
UPDATE userworkout
    SET userid = null;
