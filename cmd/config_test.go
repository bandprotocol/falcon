package cmd_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/internal/relayertest"
)

func TestConfigShow(t *testing.T) {
	t.Parallel()
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	res = sys.RunWithInput(t, "config", "show")
	require.NoError(t, res.Err)

	actual := res.Stdout.String()
	require.Equal(t, relayertest.DefaultCfgText+"\n", actual)
}

func TestConfigShowNotInit(t *testing.T) {
	t.Parallel()
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "show")
	require.ErrorContains(t, res.Err, "config does not exist:")
}

func TestConfigInitDefault(t *testing.T) {
	t.Parallel()
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

func TestConfigInitWithFileShortFlag(t *testing.T) {
	t.Parallel()
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

func TestConfigInitWithFileLongFlag(t *testing.T) {
	t.Parallel()
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

func TestConfigInitWithFileTimeString(t *testing.T) {
	t.Parallel()
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

func TestConfigInitInvalidFile(t *testing.T) {
	t.Parallel()
	sys := relayertest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	err := os.WriteFile(customCfgPath, []byte(`[band]][]]`), 0o600)
	require.NoError(t, err)

	res := sys.RunWithInput(t, "config", "init", "--file", customCfgPath)
	require.ErrorContains(t, res.Err, "error toml: expected newline")
}

func TestConfigInitNoCustomFile(t *testing.T) {
	t.Parallel()
	sys := relayertest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	res := sys.RunWithInput(t, "config", "init", "--file", customCfgPath)
	require.ErrorContains(t, res.Err, "no such file or directory")
}

func TestConfigInitAlreadyExist(t *testing.T) {
	t.Parallel()
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	res = sys.RunWithInput(t, "config", "init")
	require.ErrorContains(t, res.Err, "config already exists")
}
