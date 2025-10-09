-- +goose Up
ALTER TABLE transactions
  ADD COLUMN sender TEXT NULL;

-- Force created_at and updated_at to be stored in UTC
-- PostgreSQL TIMESTAMPTZ already stores UTC internally, but we can normalize existing data and defaults.
ALTER TABLE transactions
  ALTER COLUMN created_at SET DEFAULT (NOW() AT TIME ZONE 'UTC'),
  ALTER COLUMN updated_at SET DEFAULT (NOW() AT TIME ZONE 'UTC');

UPDATE transactions
  SET created_at = created_at AT TIME ZONE 'UTC',
      updated_at = updated_at AT TIME ZONE 'UTC';

-- Create new table for senders
CREATE TABLE senders (
  address TEXT PRIMARY KEY,
  balance NUMERIC NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC')
);

-- +goose Down
DROP TABLE IF EXISTS senders;

ALTER TABLE transactions
  DROP COLUMN IF EXISTS sender;

ALTER TABLE transactions
  ALTER COLUMN created_at SET DEFAULT NOW(),
  ALTER COLUMN updated_at SET DEFAULT NOW();
