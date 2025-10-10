-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE transactions (
  id                    INTEGER PRIMARY KEY AUTOINCREMENT,
  tx_hash               TEXT UNIQUE NOT NULL,
  tunnel_id             INTEGER NOT NULL,
  sequence              INTEGER NOT NULL,
  chain_name            TEXT NOT NULL,
  chain_type            TEXT NOT NULL CHECK (chain_type IN ('evm')),
  status                TEXT NOT NULL CHECK (status IN ('Pending','Success','Failed','Timeout')),
  sender                TEXT,
  gas_used              DECIMAL NULL,
  effective_gas_price   DECIMAL NULL,
  balance_delta         DECIMAL NULL,
  block_timestamp       DATETIME NULL,
  created_at            DATETIME NOT NULL DEFAULT (datetime('now')),
  updated_at            DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE signal_prices (
  transaction_id  INTEGER NOT NULL,
  signal_id    TEXT    NOT NULL,
  price           INTEGER NOT NULL,
  PRIMARY KEY (transaction_id, signal_id),
  FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE
);

-- +goose Down
PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS signal_prices;
DROP TABLE IF EXISTS transactions;
PRAGMA foreign_keys = ON;
