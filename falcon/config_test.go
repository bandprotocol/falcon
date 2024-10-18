package falcon_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/falcon"
	falcon_test "github.com/bandprotocol/falcon/falcon/falcontest"
)

func TestShowConfig(t *testing.T) {
	sys := falcon_test.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	res = sys.RunWithInput(t, "config", "show")
	require.NoError(t, res.Err)

	actual := res.Stdout.String()
	expect := "target_chains = []\nchecking_packet_interval = 60000000000\n\n[bandchain]\nrpc_endpoints = ['http://localhost:26657']\ntimeout = 5\n\n"

	require.Equal(t, expect, actual)
}

func TestShowEmptyConfig(t *testing.T) {
	sys := falcon_test.NewSystem(t)

	res := sys.RunWithInput(t, "config", "show")
	require.ErrorContains(t, res.Err, "config does not exist:")
}

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
