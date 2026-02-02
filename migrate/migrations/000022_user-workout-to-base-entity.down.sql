ALTER TABLE user_workout RENAME TO userworkout

ALTER TABLE user_workout DROP COLUMN updated_at;

ALTER TABLE user_workout RENAME COLUMN user_id TO userid;
ALTER TABLE user_workout RENAME COLUMN date_start TO datestart;
ALTER TABLE user_workout RENAME COLUMN date_end TO dateend;
ALTER TABLE user_workout RENAME COLUMN created_at TO createdat
