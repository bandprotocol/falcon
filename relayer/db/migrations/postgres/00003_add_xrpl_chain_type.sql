-- +goose Up
ALTER TYPE chain_type ADD VALUE IF NOT EXISTS 'xrpl';

-- +goose Down
-- NOTE: PostgreSQL does not support removing values from an ENUM type.
-- To "undo" this, would typically have to drop and recreate the type, 
-- which is dangerous if data already exists.