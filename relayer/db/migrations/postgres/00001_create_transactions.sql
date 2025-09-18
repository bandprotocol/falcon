-- +goose Up
CREATE TYPE tx_status AS ENUM ('Pending', 'Success', 'Failed', 'Timeout');
CREATE TYPE chain_type AS ENUM ('evm');

CREATE TABLE transactions (
  id                    BIGSERIAL PRIMARY KEY,
  tx_hash               TEXT UNIQUE NOT NULL,
  tunnel_id             BIGINT NOT NULL,
  sequence              BIGINT NOT NULL,
  chain_name            TEXT NOT NULL,
  chain_type            chain_type NOT NULL,
  status                tx_status NOT NULL,
  gas_used              NUMERIC NULL,
  effective_gas_price   NUMERIC NULL,
  balance_delta         NUMERIC NULL,
  block_timestamp       TIMESTAMPTZ NULL,
  created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE signal_prices (
  transaction_id  BIGINT NOT NULL,
  signal_id       TEXT   NOT NULL,
  price           BIGINT NOT NULL,
  PRIMARY KEY (transaction_id, signal_id),
  CONSTRAINT fk_signal_prices_tx
    FOREIGN KEY (transaction_id)
    REFERENCES transactions(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE signal_prices;
DROP TABLE transactions;
DROP TYPE  tx_status;
DROP TYPE  chain_type;
