package cmd_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/cmd"
	"github.com/bandprotocol/falcon/internal/relayertest"
)

func TestShowConfigCmd(t *testing.T) {
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	res = sys.RunWithInput(t, "config", "show")
	require.NoError(t, res.Err)

	actual := res.Stdout.String()
	require.Equal(t, relayertest.DefaultCfgText+"\n", actual)
}

func TestShowConfigCmdNotInit(t *testing.T) {
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "show")
	require.ErrorContains(t, res.Err, "config does not exist:")
}

func TestInitCmdDefault(t *testing.T) {
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	// read the file
	actualBytes, err := os.ReadFile(cfgPath)
	require.NoError(t, err)

	require.Equal(t, relayertest.DefaultCfgText, string(actualBytes))
}

func TestInitCmdWithFileShortFlag(t *testing.T) {
	sys := relayertest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	err := os.WriteFile(customCfgPath, []byte(relayertest.CustomCfgText), 0o600)
	require.NoError(t, err)

	res := sys.RunWithInput(t, "config", "init", "-f", customCfgPath)
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	// read the file
	actualBytes, err := os.ReadFile(cfgPath)
	require.NoError(t, err)

	require.Equal(t, relayertest.CustomCfgText, string(actualBytes))
}

func TestInitCmdWithFileLongFlag(t *testing.T) {
	sys := relayertest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	err := os.WriteFile(customCfgPath, []byte(relayertest.CustomCfgText), 0o600)
	require.NoError(t, err)

	res := sys.RunWithInput(t, "config", "init", "--file", customCfgPath)
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	actualBytes, err := os.ReadFile(cfgPath)
	require.NoError(t, err)

	require.Equal(t, relayertest.CustomCfgText, string(actualBytes))
}

func TestInitCmdWithFileTimeString(t *testing.T) {
	sys := relayertest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	err := os.WriteFile(customCfgPath, []byte(relayertest.CustomCfgTextWithTimeStr), 0o600)
	require.NoError(t, err)

	res := sys.RunWithInput(t, "config", "init", "--file", customCfgPath)
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	actualBytes, err := os.ReadFile(cfgPath)
	require.NoError(t, err)

	require.Equal(t, relayertest.CustomCfgText, string(actualBytes))
}

func TestInitCmdInvalidFile(t *testing.T) {
	sys := relayertest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	err := os.WriteFile(customCfgPath, []byte(`[band]][]]`), 0o600)
	require.NoError(t, err)

	res := sys.RunWithInput(t, "config", "init", "--file", customCfgPath)
	require.ErrorContains(t, res.Err, "error toml: expected newline")
}

func TestInitCmdNoCustomFile(t *testing.T) {
	sys := relayertest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	res := sys.RunWithInput(t, "config", "init", "--file", customCfgPath)
	require.ErrorContains(t, res.Err, "no such file or directory")
}

func TestInitCmdAlreadyExist(t *testing.T) {
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	res = sys.RunWithInput(t, "config", "init")
	require.Error(t, res.Err, cmd.ErrConfigNotExist(sys.HomeDir))
}
