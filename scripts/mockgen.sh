#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=relayer/chains/config.go -package mocks -destination internal/relayertest/mocks/chain_provider_config.go
$mockgen_cmd -source=relayer/chains/provider.go -package mocks -destination internal/relayertest/mocks/chain_provider.go
