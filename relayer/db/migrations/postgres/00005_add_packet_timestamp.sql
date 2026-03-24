-- +goose Up
ALTER TABLE transactions ADD COLUMN packet_timestamp TIMESTAMPTZ NULL;

-- +goose Down
ALTER TABLE transactions DROP COLUMN packet_timestamp;
