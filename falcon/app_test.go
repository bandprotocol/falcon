package falcon_test

import (
	"os"
	"path"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/falcon"
)

func TestInitConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Test InitConfig
	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	err := app.InitConfigFile(tmpDir)
	require.NoError(t, err)

	require.FileExists(t, path.Join(tmpDir, "config", "config.toml"))

	// read the file
	b, err := os.ReadFile(path.Join(tmpDir, "config", "config.toml"))
	require.NoError(t, err)

	// unmarshal data
	actual := falcon.Config{}
	err = toml.Unmarshal(b, &actual)
	require.NoError(t, err)

	expect := falcon.DefaultConfig()

	require.Equal(t, expect, actual)
}

func TestInitExistingConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Test InitConfig
	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	err := app.InitConfigFile(tmpDir)
	require.NoError(t, err)

	err = app.InitConfigFile(tmpDir)
	require.ErrorContains(t, err, "config already exists:")
}
