package falcon_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/falcon"
	"github.com/bandprotocol/falcon/falcon/band"
)

func TestInitConfig(t *testing.T) {
	tmpDir := t.TempDir()
	customCfgPath := ""

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	err := app.InitConfigFile(tmpDir, customCfgPath)
	require.NoError(t, err)

	cfgPath := path.Join(tmpDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	// read the file
	actualByte, err := os.ReadFile(cfgPath)
	require.NoError(t, err)

	// marshal default config
	expect := falcon.DefaultConfig()
	expectBytes, err := toml.Marshal(expect)
	require.NoError(t, err)

	require.Equal(t, string(expectBytes), string(actualByte))
}

func TestInitExistingConfig(t *testing.T) {
	tmpDir := t.TempDir()
	customCfgPath := ""

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	err := app.InitConfigFile(tmpDir, customCfgPath)
	require.NoError(t, err)

	// second time should fail
	err = app.InitConfigFile(tmpDir, customCfgPath)
	require.ErrorContains(t, err, "config already exists:")
}

func TestInitCustomConfig(t *testing.T) {
	tmpDir := t.TempDir()
	customCfgPath := path.Join(tmpDir, "custom.toml")

	// Create custom config file
	cfg := `
		[target_chains]

		[global]
		checking_packet_interval = 60000000000
	
		[bandchain]
		rpc_endpoints = ['http://localhost:26659']
		timeout = 50
	`
	// write file
	err := os.WriteFile(customCfgPath, []byte(cfg), 0o600)
	require.NoError(t, err)

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	err = app.InitConfigFile(tmpDir, customCfgPath)
	require.NoError(t, err)

	require.FileExists(t, path.Join(tmpDir, "config", "config.toml"))

	// read the file
	b, err := os.ReadFile(path.Join(tmpDir, "config", "config.toml"))
	require.NoError(t, err)

	// unmarshal data
	actual := falcon.Config{}
	err = toml.Unmarshal(b, &actual)
	require.NoError(t, err)

	expect := falcon.Config{
		BandChain: band.Config{
			RpcEndpoints: []string{"http://localhost:26659"},
			Timeout:      50,
		},
		TargetChains: nil,
		Global: falcon.GlobalConfig{
			CheckingPacketInterval: time.Minute,
		},
	}

	require.Equal(t, expect, actual)
}
