package falcontest

import (
	"github.com/bandprotocol/falcon/falcon"
	"github.com/bandprotocol/falcon/falcon/band"
	"github.com/bandprotocol/falcon/falcon/chains"
	"github.com/bandprotocol/falcon/falcon/chains/evm"
)

const DefaultCfgText = `[global]
checking_packet_interval = 60000000000

[bandchain]
rpc_endpoints = ['http://localhost:26657']
timeout = 5

[target_chains]
`

const CustomCfgText = `[global]
checking_packet_interval = 1

[bandchain]
rpc_endpoints = ['http://localhost:26657', 'http://localhost:26658']
timeout = 0

[target_chains]
[target_chains.testnet]
rpc_endpoints = ['http://localhost:8545']
chain_type = 'evm'
max_retry = 3
chain_id = 31337
tunnel_router_address = '0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9'
`

var CustomCfg = falcon.Config{
	Global: falcon.GlobalConfig{CheckingPacketInterval: 1},
	BandChain: band.Config{
		RpcEndpoints: []string{"http://localhost:26657", "http://localhost:26658"},
		Timeout:      0,
	},
	TargetChains: chains.ChainProviderConfigs{
		"testnet": &evm.EVMChainProviderConfig{
			BaseChainProviderConfig: chains.BaseChainProviderConfig{
				RpcEndpoints:        []string{"http://localhost:8545"},
				ChainType:           chains.ChainTypeEVM,
				MaxRetry:            3,
				ChainID:             31337,
				TunnelRouterAddress: "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9",
			},
		},
	},
}
