package cmd_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/internal/falcontest"
)

func TestShowConfigCmd(t *testing.T) {
	sys := falcontest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	res = sys.RunWithInput(t, "config", "show")
	require.NoError(t, res.Err)

	actual := res.Stdout.String()
	require.Equal(t, falcontest.DefaultCfgText+"\n", actual)
}

func TestShowConfigCmdNotInit(t *testing.T) {
	sys := falcontest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "show")
	require.ErrorContains(t, res.Err, "config does not exist:")
}

func TestInitCmdDefault(t *testing.T) {
	sys := falcontest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	// read the file
	actualBytes, err := os.ReadFile(cfgPath)
	require.NoError(t, err)

	require.Equal(t, falcontest.DefaultCfgText, string(actualBytes))
}

func TestInitCmdWithFileShortFlag(t *testing.T) {
	sys := falcontest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	err := os.WriteFile(customCfgPath, []byte(falcontest.CustomCfgText), 0o600)
	require.NoError(t, err)

	res := sys.RunWithInput(t, "config", "init", "-f", customCfgPath)
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	// read the file
	actualBytes, err := os.ReadFile(cfgPath)
	require.NoError(t, err)

	require.Equal(t, falcontest.CustomCfgText, string(actualBytes))
}

func TestInitCmdWithFileLongFlag(t *testing.T) {
	sys := falcontest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	err := os.WriteFile(customCfgPath, []byte(falcontest.CustomCfgText), 0o600)
	require.NoError(t, err)

	res := sys.RunWithInput(t, "config", "init", "--file", customCfgPath)
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	actualBytes, err := os.ReadFile(cfgPath)
	require.NoError(t, err)

	require.Equal(t, falcontest.CustomCfgText, string(actualBytes))
}

func TestInitCmdNoCustomFile(t *testing.T) {
	sys := falcontest.NewSystem(t)

	customCfgPath := path.Join(sys.HomeDir, "custom.toml")
	res := sys.RunWithInput(t, "config", "init", "--file", customCfgPath)
	require.ErrorContains(t, res.Err, "no such file or directory")
}

func TestInitCmdAlreadyExist(t *testing.T) {
	sys := falcontest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	cfgPath := path.Join(sys.HomeDir, "config", "config.toml")
	require.FileExists(t, cfgPath)

	res = sys.RunWithInput(t, "config", "init")
	require.ErrorContains(t, res.Err, "config already exists")
}