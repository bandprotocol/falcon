-- +goose Up
-- Add 'sender' column to transactions (SQLite cannot use ALTER COLUMN)
ALTER TABLE transactions ADD COLUMN sender TEXT;

-- Note: SQLite doesn't have TIMESTAMPTZ; it stores times as TEXT or REAL.
-- We assume UTC timestamps are inserted by the application (Go, etc.)
-- so no need to alter existing data.

-- Create new table for senders
CREATE TABLE senders (
  address TEXT PRIMARY KEY,
  balance NUMERIC NOT NULL,
  created_at DATETIME NOT NULL DEFAULT (datetime('now')), -- UTC by default
  updated_at DATETIME NOT NULL DEFAULT (datetime('now'))  -- UTC by default
);

-- +goose Down
DROP TABLE IF EXISTS senders;

-- SQLite cannot drop columns directly, so we need to recreate the table if we want to fully revert.
-- For rollback simplicity, we’ll just leave the 'sender' column unused.
-- (If a full schema revert is required, a manual table rebuild is necessary.)
