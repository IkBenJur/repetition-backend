UPDATE users SET active_userworkout_id = (
    SELECT id FROM userworkout WHERE name = 'Upper body day [MIGRATION-default-test-data]' LIMIT 1
) WHERE id = 1
