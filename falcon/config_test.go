package falcon_test

import (
	"path"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/falcon"
	"github.com/bandprotocol/falcon/falcon/band"
	"github.com/bandprotocol/falcon/falcon/chains"
	"github.com/bandprotocol/falcon/falcon/chains/evm"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	customConfigPath := ""
	cfgPath := path.Join(tmpDir, "config", "config.toml")

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	// Prepare config before test
	err := app.InitConfigFile(tmpDir, customConfigPath)
	require.NoError(t, err)

	actual, err := falcon.LoadConfig(cfgPath)
	require.NoError(t, err)
	expect := falcon.DefaultConfig()
	require.Equal(t, expect, actual)
}

func TestUnmarshalConfig(t *testing.T) {
	// create new toml config file
	cfgText := `
		[global]
		checking_packet_interval = 60000000000
		target_chains = []

		[bandchain]
		rpc_endpoints = ['http://localhost:26657']
		timeout = 7

		[target_chains]
		[target_chains.ethereum]
		chain_type = "evm"
		rpc_endpoints = ['http://localhost:26657']
	`

	// unmarshall them with Config into struct
	cfgWrapper := &falcon.ConfigInputWrapper{}
	err := toml.Unmarshal([]byte(cfgText), cfgWrapper)
	require.NoError(t, err)
	require.Equal(t, falcon.ConfigInputWrapper{
		Global: falcon.GlobalConfig{CheckingPacketInterval: time.Minute},
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26657"},
			Timeout:      7,
		},
		TargetChains: map[string]falcon.TOMLWrapper{
			"ethereum": {
				"chain_type":    "evm",
				"rpc_endpoints": []interface{}{"http://localhost:26657"},
			},
		},
	}, *cfgWrapper)

	cfg, err := falcon.ParseConfig(cfgWrapper)
	require.NoError(t, err)

	require.Equal(t, &falcon.Config{
		Global: falcon.GlobalConfig{CheckingPacketInterval: time.Minute},
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26657"},
			Timeout:      7,
		},
		TargetChains: map[string]chains.ChainProviderConfig{
			"ethereum": &evm.EVMProviderConfig{
				RpcEndpoints: []string{"http://localhost:26657"},
				ChainType:    chains.ChainTypeEVM,
			},
		},
	}, cfg)
}

func TestMarshalConfig(t *testing.T) {
	cfg := &falcon.Config{
		Global: falcon.GlobalConfig{CheckingPacketInterval: time.Minute},
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26657"},
			Timeout:      7,
		},
		TargetChains: map[string]chains.ChainProviderConfig{
			"ethereum": &evm.EVMProviderConfig{
				RpcEndpoints: []string{"http://localhost:26657"},
				ChainType:    chains.ChainTypeEVM,
			},
		},
	}

	b, err := toml.Marshal(cfg)
	expect := `[global]
checking_packet_interval = 60000000000

[bandchain]
rpc_endpoints = ['http://localhost:26657']
timeout = 7

[target_chains]
[target_chains.ethereum]
rpc_endpoints = ['http://localhost:26657']
chain_type = 'evm'
`
	require.NoError(t, err)
	require.Equal(t, expect, string(b))
}
