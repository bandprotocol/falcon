-- +goose Up
CREATE INDEX idx_transactions_chain_name_created_at ON transactions (chain_name, created_at DESC);

-- +goose Down
DROP INDEX idx_transactions_chain_name_created_at;
