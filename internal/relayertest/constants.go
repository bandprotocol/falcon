package relayertest

import (
	_ "embed"
	"time"

	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
	"github.com/bandprotocol/falcon/relayer/config"
)

//go:embed testdata/default_config.toml
var DefaultCfgText string

//go:embed testdata/custom_config.toml
var CustomCfgText string

//go:embed testdata/custom_config_with_time_str.toml
var CustomCfgTextWithTimeStr string

var CustomCfg = config.Config{
	Global: config.GlobalConfig{
		CheckingPacketInterval: 1 * time.Minute,
		SyncTunnelsInterval:    5 * time.Minute,
		PenaltySkipRounds:      3,
		LogLevel:               "info",
	},
	BandChain: band.Config{
		RpcEndpoints:               []string{"http://localhost:26657", "http://localhost:26658"},
		Timeout:                    3 * time.Second,
		LivelinessCheckingInterval: 5 * time.Minute,
	},
	TargetChains: config.ChainProviderConfigs{
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
			LivelinessCheckingInterval: 5 * time.Minute,
			GasType:                    evm.GasTypeEIP1559,
			GasMultiplier:              1.1,
		},
	},
}

//go:embed testdata/chain_config.toml
var ChainCfgText string

//go:embed testdata/default_with_chain_config.toml
var DefaultCfgTextWithChainCfg string

//go:embed testdata/chain_config_invalid_chain_type.toml
var ChainCfgInvalidChainTypeText string
