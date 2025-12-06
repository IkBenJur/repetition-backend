DO $$
DECLARE   -- Variables to hold generated IDs
    v_userworkout_id INT;
    v_userworkout_exercise_id INT;
BEGIN

    -- Create new workout 
    INSERT INTO userworkout (name, datestart, userid)
    VALUES ('Upper body day [MIGRATION-default-test-data]', current_timestamp, 1)
    RETURNING id INTO v_userworkout_id;

        -- Add exercise to workout
        INSERT INTO userworkoutexercise (userworkoutid, exerciseid)
        VALUES (v_userworkout_id, 2)
        RETURNING id INTO v_userworkout_exercise_id;

            -- Add sets
            INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight)
            VALUES 
                (v_userworkout_exercise_id, 8, 100),
                (v_userworkout_exercise_id, 6, 100),
                (v_userworkout_exercise_id, 10, 90);

        -- Another exercise
        INSERT INTO userworkoutexercise (userworkoutid, exerciseid)
        VALUES (v_userworkout_id, 4)
        RETURNING id INTO v_userworkout_exercise_id;

            -- Add sets
            INSERT INTO userworkoutexerciseset (userworkoutexerciseid, reps, weight)
            VALUES 
                (v_userworkout_exercise_id, 12, 100),
                (v_userworkout_exercise_id, 10, 100),
                (v_userworkout_exercise_id, 10, 100);
 
END $$;