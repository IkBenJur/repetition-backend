ALTER TABLE user_table RENAME TO users;

ALTER TABLE users DROP COLUMN updated_at;
ALTER TABLE users RENAME COLUMN created_at TO createdat;
