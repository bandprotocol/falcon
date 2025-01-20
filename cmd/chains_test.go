package cmd_test

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/internal/relayertest"
)

func TestChainsListEmpty(t *testing.T) {
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	res = sys.RunWithInput(t, "chains", "list")
	require.Empty(t, res.Stdout.String())
}

func TestChainsAdd(t *testing.T) {
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	chainCfgPath := path.Join(sys.HomeDir, "chain_config.toml")
	err := os.WriteFile(chainCfgPath, []byte(relayertest.ChainCfgText), 0o600)
	require.NoError(t, err)

	require.FileExists(t, chainCfgPath)

	// Add chain
	res = sys.RunWithInput(t, "chains", "add", "testnet", chainCfgPath)
	require.Empty(t, res.Stdout.String())
	require.Empty(t, res.Stderr.String())

	// Add another chain
	res = sys.RunWithInput(t, "ch", "a", "testnet2", chainCfgPath)
	require.Empty(t, res.Stdout.String())
	require.Empty(t, res.Stderr.String())

	// Add existing chain
	res = sys.RunWithInput(t, "ch", "a", "testnet", chainCfgPath)
	require.Empty(t, res.Stdout.String())
	require.Error(t, res.Err, "existing chain name")

	// List chains to check
	res = sys.RunWithInput(t, "chains", "list")
	require.Regexp(t, regexp.MustCompile(`\d+: ([\w-]+) -> type\((\w+)\)`), res.Stdout.String())
	require.Empty(t, res.Stderr.String())
}

func TestChainsDelete(t *testing.T) {
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	chainCfgPath := path.Join(sys.HomeDir, "chain_config.toml")
	err := os.WriteFile(chainCfgPath, []byte(relayertest.ChainCfgText), 0o600)
	require.NoError(t, err)

	require.FileExists(t, chainCfgPath)

	// Add chain
	res = sys.RunWithInput(t, "chains", "add", "testnet", chainCfgPath)
	require.Empty(t, res.Stdout.String())
	require.Empty(t, res.Stderr.String())

	// Add another chain
	res = sys.RunWithInput(t, "chains", "add", "testnet2", chainCfgPath)
	require.Empty(t, res.Stdout.String())
	require.Empty(t, res.Stderr.String())

	// List chains
	res = sys.RunWithInput(t, "chains", "list")
	require.Regexp(t, regexp.MustCompile(`\d+: ([\w-]+) -> type\((\w+)\)`), res.Stdout.String())
	require.Empty(t, res.Stderr.String())

	// Delete chain
	res = sys.RunWithInput(t, "chains", "delete", "testnet")
	require.Empty(t, res.Stdout.String())
	require.Empty(t, res.Stderr.String())

	res = sys.RunWithInput(t, "ch", "d", "testnet2")
	require.Empty(t, res.Stdout.String())
	require.Empty(t, res.Stderr.String())

	// List chain with shorthand command
	res = sys.RunWithInput(t, "ch", "l")
	require.Empty(t, res.Stdout.String())
}

func TestChainsShow(t *testing.T) {
	sys := relayertest.NewSystem(t)

	res := sys.RunWithInput(t, "config", "init")
	require.NoError(t, res.Err)

	chainCfgPath := path.Join(sys.HomeDir, "chain_config.toml")
	err := os.WriteFile(chainCfgPath, []byte(relayertest.ChainCfgText), 0o600)
	require.NoError(t, err)

	require.FileExists(t, chainCfgPath)

	// Add chain
	res = sys.RunWithInput(t, "chains", "add", "testnet", chainCfgPath)
	require.Empty(t, res.Stdout.String())
	require.Empty(t, res.Stderr.String())

	// Show chain configuration
	res = sys.RunWithInput(t, "chains", "show", "testnet")

	var expectedChainCfg map[string]interface{}
	err = toml.Unmarshal(res.Stdout.Bytes(), &expectedChainCfg)
	require.NoError(t, err)

	var actualChainCfg map[string]interface{}
	err = toml.Unmarshal([]byte(relayertest.ChainCfgText), &actualChainCfg)
	require.NoError(t, err)

	require.Equal(t, expectedChainCfg, actualChainCfg)
}
