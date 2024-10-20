#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=falcon/chains/config.go -package mocks -destination internal/falcontest/mocks/chain_provider_config.go
$mockgen_cmd -source=falcon/chains/provider.go -package mocks -destination internal/falcontest/mocks/chain_provider.go
