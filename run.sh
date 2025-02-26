#!/bin/bash

# Copy config from config map mount path
falcon config init
cp /config/config.toml /app/.falcon/config/config.toml

# Add keys to key ring
falcon keys add holesky-testnet testkey-1 --private-key $ETH_PRIV_KEY

# Start Service
falcon start
