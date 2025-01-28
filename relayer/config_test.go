package relayer_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/internal/relayertest"
	"github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := path.Join(tmpDir, "config", "config.toml")

	testcases := []struct {
		name        string
		preProcess  func(t *testing.T)
		postProcess func(t *testing.T)
		out         *relayer.Config
		err         error
	}{
		{
			name: "read default config",
			preProcess: func(t *testing.T) {
				app := relayer.NewApp(nil, tmpDir, false, nil)
				err := app.InitConfigFile(tmpDir, "")
				require.NoError(t, err)
			},
			out: relayer.DefaultConfig(),
			postProcess: func(t *testing.T) {
				err := os.Remove(cfgPath)
				require.NoError(t, err)
			},
		},
		{
			name: "no config file",
			err:  fmt.Errorf("no such file or directory"),
		},
		{
			name: "invalid config file; invalid chain type",
			preProcess: func(t *testing.T) {
				// create new toml config file
				cfgText := `[target_chains.testnet]
			chain_type = 'evms'
			`

				err := os.WriteFile(cfgPath, []byte(cfgText), 0o600)
				require.NoError(t, err)
			},
			err: relayer.ErrUnsupportedChainType("evms"),
			postProcess: func(t *testing.T) {
				err := os.Remove(cfgPath)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.preProcess != nil {
				tc.preProcess(t)
			}

			if tc.postProcess != nil {
				defer tc.postProcess(t)
			}

			actual, err := relayer.LoadConfig(cfgPath)
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
		in   relayer.ChainProviderConfigWrapper
		out  chains.ChainProviderConfig
		err  error
	}{
		{
			name: "valid evm chain",
			in: relayer.ChainProviderConfigWrapper{
				"chain_type": "evm",
				"endpoints":  []string{"http://localhost:8545"},
			},
			out: &evm.EVMChainProviderConfig{
				BaseChainProviderConfig: chains.BaseChainProviderConfig{
					Endpoints: []string{"http://localhost:8545"},
					ChainType: chains.ChainTypeEVM,
				},
			},
		},
		{
			name: "chain type not found",
			in: relayer.ChainProviderConfigWrapper{
				"chain_type": "evms",
				"endpoints":  []string{"http://localhost:8545"},
			},
			err: fmt.Errorf("unsupported chain type: evms"),
		},
		{
			name: "missing chain type",
			in: relayer.ChainProviderConfigWrapper{
				"endpoints": []string{"http://localhost:8545"},
			},
			err: fmt.Errorf("chain_type is required"),
		},
		{
			name: "chain type not string",
			in: relayer.ChainProviderConfigWrapper{
				"chain_type": []string{"evm"},
				"endpoints":  []string{"http://localhost:8545"},
			},
			err: fmt.Errorf("chain_type is required"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := relayer.ParseChainProviderConfig(tc.in)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.out, actual)
			}
		})
	}
}

func TestParseConfigInvalidChainProviderConfig(t *testing.T) {
	w := &relayer.ConfigInputWrapper{
		Global: relayer.GlobalConfig{CheckingPacketInterval: 1},
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26657", "http://localhost:26658"},
			Timeout:      0,
		},
		TargetChains: map[string]relayer.ChainProviderConfigWrapper{
			"testnet": {
				"chain_type": "evms",
			},
		},
	}

	_, err := relayer.ParseConfig(w)
	require.Error(t, err, relayer.ErrUnsupportedChainType("evms"))
}

func TestUnmarshalConfig(t *testing.T) {
	// create new toml config file
	cfgText := relayertest.CustomCfgText

	// unmarshal them with Config into struct
	var cfgWrapper relayer.ConfigInputWrapper
	err := relayer.DecodeConfigInputWrapperTOML([]byte(cfgText), &cfgWrapper)
	require.NoError(t, err)

	cfg, err := relayer.ParseConfig(&cfgWrapper)
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
	actual, err := relayer.LoadChainConfig(cfgPath)
	require.NoError(t, err)

	expect := relayertest.CustomCfg.TargetChains[chainName]

	require.Equal(t, expect, actual)
}
