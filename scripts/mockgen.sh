#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=relayer/band/query.go -package mocks -destination internal/relayertest/mocks/band_chain_query.go
$mockgen_cmd -source=relayer/chains/config.go -package mocks -destination internal/relayertest/mocks/chain_provider_config.go
$mockgen_cmd -source=relayer/chains/provider.go -package mocks -destination internal/relayertest/mocks/chain_provider.go
$mockgen_cmd -source=relayer/chains/evm/client.go -mock_names Client=MockEVMClient -package mocks -destination internal/relayertest/mocks/chain_evm_client.go
$mockgen_cmd -source=relayer/band/client.go -package mocks -destination internal/relayertest/mocks/band_client.go
$mockgen_cmd -source=relayer/wallet/wallet.go -package mocks -destination internal/relayertest/mocks/wallet.go
$mockgen_cmd -source=relayer/store/store.go -package mocks -destination internal/relayertest/mocks/store.go
$mockgen_cmd -source=proto/kms/signer_grpc.pb.go -package mocks -destination internal/relayertest/mocks/signer_grpc.go