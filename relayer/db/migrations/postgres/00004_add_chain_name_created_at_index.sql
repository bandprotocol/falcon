-- +goose Up
CREATE INDEX idx_transactions_tunnel_id_block_timestamp ON transactions (tunnel_id, block_timestamp DESC);

-- +goose Down
DROP INDEX idx_transactions_tunnel_id_block_timestamp;
