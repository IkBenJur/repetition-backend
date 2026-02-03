ALTER TABLE user_workout RENAME TO userworkout;

ALTER TABLE userworkout DROP COLUMN updated_at;

ALTER TABLE userworkout RENAME COLUMN user_id TO userid;
ALTER TABLE userworkout RENAME COLUMN date_start TO datestart;
ALTER TABLE userworkout RENAME COLUMN date_end TO dateend;
ALTER TABLE userworkout RENAME COLUMN created_at TO createdat
