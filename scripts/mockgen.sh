#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=relayer/chains/config.go -package mocks -destination internal/relayertest/mocks/chain_provider_config.go
$mockgen_cmd -source=relayer/chains/provider.go -package mocks -destination internal/relayertest/mocks/chain_provider.go
$mockgen_cmd -source=relayer/band/client.go -package mocks -destination internal/relayertest/mocks/band_client.go
$mockgen_cmd -source=relayer/band/provider.go -package mocks -destination internal/relayertest/mocks/band_provider.go
