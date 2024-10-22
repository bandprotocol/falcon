package relayer_test

import (
	"os"
	"path"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/internal/relayertest"
	falcon "github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
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

func TestLoadConfigNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := path.Join(tmpDir, "config", "config.toml")

	_, err := falcon.LoadConfig(cfgPath)
	require.ErrorContains(t, err, "no such file or directory")
}

func TestLoadConfigInvalidChainProviderConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := path.Join(tmpDir, "config.toml")

	// create new toml config file
	cfgText := `[target_chains.testnet]
chain_type = 'evms'
`

	err := os.WriteFile(cfgPath, []byte(cfgText), 0o600)
	require.NoError(t, err)

	_, err = falcon.LoadConfig(cfgPath)
	require.ErrorContains(t, err, "unsupported chain type: evms")
}

func TestParseChainProviderConfigTypeEVM(t *testing.T) {
	w := falcon.TOMLWrapper{
		"chain_type": "evm",
		"endpoints":  []string{"http://localhost:8545"},
	}

	cfg, err := falcon.ParseChainProviderConfig(w)

	expect := &evm.EVMChainProviderConfig{
		BaseChainProviderConfig: chains.BaseChainProviderConfig{
			Endpoints: []string{"http://localhost:8545"},
			ChainType: chains.ChainTypeEVM,
		},
	}
	require.NoError(t, err)
	require.Equal(t, expect, cfg)
}

func TestParseChainProviderConfigTypeNotFound(t *testing.T) {
	w := falcon.TOMLWrapper{
		"chain_type": "evms",
		"endpoints":  []string{"http://localhost:8545"},
	}

	_, err := falcon.ParseChainProviderConfig(w)
	require.ErrorContains(t, err, "unsupported chain type: evms")
}

func TestParseChainProviderConfigNoChainType(t *testing.T) {
	w := falcon.TOMLWrapper{
		"endpoints": []string{"http://localhost:8545"},
	}

	_, err := falcon.ParseChainProviderConfig(w)
	require.ErrorContains(t, err, "chain_type is required")
}

func TestParseConfigInvalidChainProviderConfig(t *testing.T) {
	w := &falcon.ConfigInputWrapper{
		Global: falcon.GlobalConfig{CheckingPacketInterval: 1},
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26657", "http://localhost:26658"},
			Timeout:      0,
		},
		TargetChains: map[string]falcon.TOMLWrapper{
			"testnet": {
				"chain_type": "evms",
			},
		},
	}

	_, err := falcon.ParseConfig(w)
	require.ErrorContains(t, err, "unsupported chain type: evms")
}

func TestUnmarshalConfig(t *testing.T) {
	// create new toml config file
	cfgText := relayertest.CustomCfgText

	// unmarshal them with Config into struct
	cfgWrapper := &falcon.ConfigInputWrapper{}
	err := toml.Unmarshal([]byte(cfgText), cfgWrapper)
	require.NoError(t, err)

	cfg, err := falcon.ParseConfig(cfgWrapper)
	require.NoError(t, err)

	require.Equal(t, &relayertest.CustomCfg, cfg)
}

func TestMarshalConfig(t *testing.T) {
	b, err := toml.Marshal(relayertest.CustomCfg)
	require.NoError(t, err)
	require.Equal(t, relayertest.CustomCfgText, string(b))
}
