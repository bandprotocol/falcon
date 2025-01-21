#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=relayer/chains/config.go -package mocks -destination internal/relayertest/mocks/chain_provider_config.go
$mockgen_cmd -source=relayer/chains/provider.go -package mocks -destination internal/relayertest/mocks/chain_provider.go
$mockgen_cmd -source=relayer/chains/evm/client.go -mock_names Client=MockEVMClient -package mocks -destination internal/relayertest/mocks/chain_evm_client.go
$mockgen_cmd -source=relayer/band/client.go -package mocks -destination internal/relayertest/mocks/band_client.go
$mockgen_cmd -package mocks -mock_names QueryClient=MockTunnelQueryClient -destination internal/relayertest/mocks/tunnel_query_client.go github.com/bandprotocol/chain/v3/x/tunnel/types QueryClient
$mockgen_cmd -package mocks -mock_names QueryClient=MockBandtssQueryClient -destination internal/relayertest/mocks/bandtss_query_client.go github.com/bandprotocol/chain/v3/x/bandtss/types QueryClient

