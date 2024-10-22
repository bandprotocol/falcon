package falcontest

import (
	_ "embed"
	"time"

	"github.com/bandprotocol/falcon/falcon"
	"github.com/bandprotocol/falcon/falcon/band"
	"github.com/bandprotocol/falcon/falcon/chains"
	"github.com/bandprotocol/falcon/falcon/chains/evm"
)

//go:embed testdata/default_config.toml
var DefaultCfgText string

//go:embed testdata/custom_config.toml
var CustomCfgText string

var CustomCfg = falcon.Config{
	Global: falcon.GlobalConfig{CheckingPacketInterval: 1},
	BandChain: band.Config{
		RpcEndpoints: []string{"http://localhost:26657", "http://localhost:26658"},
		Timeout:      0,
	},
	TargetChains: chains.ChainProviderConfigs{
		"testnet": &evm.EVMChainProviderConfig{
			BaseChainProviderConfig: chains.BaseChainProviderConfig{
				Endpoints:           []string{"http://localhost:8545"},
				ChainType:           chains.ChainTypeEVM,
				MaxRetry:            3,
				ChainID:             31337,
				TunnelRouterAddress: "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9",
				QueryTimeout:        time.Second * 3,
			},
		},
	},
}
