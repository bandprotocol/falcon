package relayertest

import (
	_ "embed"
	"time"

	falcon "github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
)

//go:embed testdata/default_config.toml
var DefaultCfgText string

//go:embed testdata/custom_config.toml
var CustomCfgText string

//go:embed testdata/custom_config_with_time_str.toml
var CustomCfgTextWithTimeStr string

var CustomCfg = falcon.Config{
	Global: falcon.GlobalConfig{
		CheckingPacketInterval:           1 * time.Minute,
		MaxCheckingPacketPenaltyDuration: 1 * time.Hour,
		PenaltyExponentialFactor:         1.1,
		LogLevel:                         "info",
	},
	BandChain: band.Config{
		RpcEndpoints: []string{"http://localhost:26657", "http://localhost:26658"},
		Timeout:      3 * time.Second,
	},
	TargetChains: chains.ChainProviderConfigs{
		"testnet": &evm.EVMChainProviderConfig{
			BaseChainProviderConfig: chains.BaseChainProviderConfig{
				Endpoints:           []string{"http://localhost:8545"},
				ChainType:           chains.ChainTypeEVM,
				MaxRetry:            3,
				ChainID:             31337,
				TunnelRouterAddress: "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9",
				QueryTimeout:        3 * time.Second,
				ExecuteTimeout:      3 * time.Second,
			},
			BlockConfirmation:          5,
			WaitingTxDuration:          time.Second * 3,
			CheckingTxInterval:         time.Second,
			LivelinessCheckingInterval: 15 * time.Minute,
			GasType:                    evm.GasTypeEIP1559,
			GasMultiplier:              1.1,
		},
	},
}
