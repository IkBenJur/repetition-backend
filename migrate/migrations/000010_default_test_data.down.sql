-- 1. Delete sets
DELETE FROM userworkoutexerciseset
WHERE userworkoutexerciseid IN (
    SELECT id 
    FROM userworkoutexercise 
    WHERE userworkoutid IN (
        SELECT id 
        FROM userworkout 
        WHERE userid = 1 AND name = 'Upper body day [MIGRATION-default-test-data]'
    )
);

-- 2. Delete exercises
DELETE FROM userworkoutexercise
WHERE userworkoutid IN (
    SELECT id 
    FROM userworkout 
    WHERE userid = 1 AND name = 'Upper body day [MIGRATION-default-test-data]'
);

-- 3. Delete workout
DELETE FROM userworkout 
WHERE userid = 1 
  AND name = 'Upper body day [MIGRATION-default-test-data]';