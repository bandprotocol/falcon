package config_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/internal/relayertest"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/config"
)

func TestParseConfig(t *testing.T) {
	testcases := []struct {
		name        string
		in          []byte
		preProcess  func(t *testing.T)
		postProcess func(t *testing.T)
		out         *config.Config
		err         error
	}{
		{
			name: "read default config",
			in:   []byte(relayertest.DefaultCfgText),
			out:  config.DefaultConfig(),
		},
		{
			name: "invalid config file; invalid chain type",
			in: []byte(`[target_chains.testnet]
			chain_type = 'evms'
			`),
			err: fmt.Errorf("unsupported chain type: evms"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := config.ParseConfig(tc.in)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.out, actual)
			}
		})
	}
}

func TestParseChainProviderConfig(t *testing.T) {
	testcases := []struct {
		name string
		in   config.ChainProviderConfigWrapper
		out  chains.ChainProviderConfig
		err  error
	}{
		{
			name: "valid evm chain",
			in: config.ChainProviderConfigWrapper{
				"chain_type": "evm",
				"endpoints":  []string{"http://localhost:8545"},
			},
			out: &evm.EVMChainProviderConfig{
				BaseChainProviderConfig: chains.BaseChainProviderConfig{
					Endpoints: []string{"http://localhost:8545"},
					ChainType: chainstypes.ChainTypeEVM,
				},
			},
		},
		{
			name: "chain type not found",
			in: config.ChainProviderConfigWrapper{
				"chain_type": "evms",
				"endpoints":  []string{"http://localhost:8545"},
			},
			err: fmt.Errorf("unsupported chain type: evms"),
		},
		{
			name: "missing chain type",
			in: config.ChainProviderConfigWrapper{
				"endpoints": []string{"http://localhost:8545"},
			},
			err: fmt.Errorf("chain_type is required"),
		},
		{
			name: "chain type not string",
			in: config.ChainProviderConfigWrapper{
				"chain_type": []string{"evm"},
				"endpoints":  []string{"http://localhost:8545"},
			},
			err: fmt.Errorf("chain_type is required"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := config.ParseChainProviderConfig(tc.in)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.out, actual)
			}
		})
	}
}

func TestParseConfigInputWrapperInvalidChainProviderConfig(t *testing.T) {
	w := &config.ConfigInputWrapper{
		Global: config.GlobalConfig{CheckingPacketInterval: 1},
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26657", "http://localhost:26658"},
			Timeout:      0,
		},
		TargetChains: map[string]config.ChainProviderConfigWrapper{
			"testnet": {
				"chain_type": "evms",
			},
		},
	}

	_, err := config.ParseConfigInputWrapper(w)
	require.ErrorContains(t, err, "unsupported chain type: evms")
}

func TestParseConfigInputWrapper(t *testing.T) {
	// create new toml config file
	cfgText := relayertest.CustomCfgText

	// unmarshal them with Config into struct
	var cfgWrapper config.ConfigInputWrapper
	err := config.DecodeConfigInputWrapperTOML([]byte(cfgText), &cfgWrapper)
	require.NoError(t, err)

	cfg, err := config.ParseConfigInputWrapper(&cfgWrapper)
	require.NoError(t, err)

	require.Equal(t, &relayertest.CustomCfg, cfg)
}

func TestMarshalConfig(t *testing.T) {
	b, err := toml.Marshal(relayertest.CustomCfg)
	require.NoError(t, err)
	require.Equal(t, relayertest.CustomCfgText, string(b))
}

func TestLoadChainConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := path.Join(tmpDir, "chain_config.toml")
	chainName := "testnet"

	// write config file
	err := os.WriteFile(cfgPath, []byte(relayertest.ChainCfgText), 0o600)
	require.NoError(t, err)

	// load chain config
	actual, err := config.LoadChainConfig(cfgPath)
	require.NoError(t, err)

	expect := relayertest.CustomCfg.TargetChains[chainName]

	require.Equal(t, expect, actual)
}
