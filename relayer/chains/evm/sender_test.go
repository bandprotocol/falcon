package evm_test

import (
	"testing"

	"os"
	"path"

	"github.com/bandprotocol/falcon/relayer/chains/evm"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
)

func TestLoadKeyInfo(t *testing.T) {
	tmpDir := t.TempDir()
	chainName := "testnet"

	// write mock keyInfo at keyInfo's path
	keyInfo := make(evm.KeyInfo)
	keyInfo["key1"] = ""
	keyInfo["key2"] = ""
	b, err := toml.Marshal(&keyInfo)

	keyInfoDir := path.Join(tmpDir, "keys", chainName, "info")
	keyInfoPath := path.Join(keyInfoDir, "info.toml")
	// Create the info folder if doesn't exist
	err = os.Mkdir(keyInfoDir, os.ModePerm)
	require.NoError(t, err)
	// Create the file and write the default config to the given location.
	f, err := os.Create(keyInfoPath)
	require.NoError(t, err)
	defer f.Close()

	_, err = f.Write(b)
	require.NoError(t, err)

	// load keyInfo
	actual, err := evm.LoadKeyInfo(tmpDir, chainName)
	require.NoError(t, err)

	require.Equal(t, keyInfo, actual)
}
