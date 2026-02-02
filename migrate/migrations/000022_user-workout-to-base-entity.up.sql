ALTER TABLE userworkout RENAME TO user_workout;

ALTER TABLE user_workout ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE user_workout RENAME COLUMN userid TO user_id;
ALTER TABLE user_workout RENAME COLUMN datestart TO date_start;
ALTER TABLE user_workout RENAME COLUMN dateend TO date_end;
ALTER TABLE user_workout RENAME COLUMN createdat TO created_at;
