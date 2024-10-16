package falcon_test

import (
	"fmt"
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
	customCfgPath := ""

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	err := app.InitConfigFile(tmpDir, customCfgPath)
	require.NoError(t, err)

	// second time should fail
	err = app.InitConfigFile(tmpDir, customCfgPath)
	require.ErrorContains(t, err, "config already exists:")
}

func TestShowEmptyConfig(t *testing.T) {
	tmpDir := t.TempDir()

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	// should not have the file
	require.NoFileExists(t, path.Join(tmpDir, "config", "config.toml"))

	_, err := app.GetConfigFile(tmpDir)
	require.ErrorContains(t, err, "config does not exist:")
}

func TestShowConfig(t *testing.T) {
	tmpDir := t.TempDir()
	customCfgPath := ""

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	// create config file
	err := app.InitConfigFile(tmpDir, customCfgPath)
	require.NoError(t, err)

	require.FileExists(t, path.Join(tmpDir, "config", "config.toml"))

	// assign config
	cfg := falcon.DefaultConfig()
	app.Config = &(cfg)

	// read config
	cfgContent, err := app.GetConfigFile(tmpDir)
	require.NoError(t, err)

	require.Contains(t, cfgContent, "target_chains = []")
	require.Contains(t, cfgContent, "checking_packet_interval = 60000000000")
	require.Contains(t, cfgContent, "[bandchain]")
	require.Contains(t, cfgContent, "rpc_endpoints = ['http://localhost:26657']")
	require.Contains(t, cfgContent, "timeout = 5")
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// custom config path
	customCfgPath := path.Join(tmpDir, "custom.toml")

	// create new toml config file
	cfg := `
		target_chains = []
		checking_packet_interval = 60000000000
	
		[bandchain]
		rpc_endpoints = ['http://localhost:26657']
		timeout = 7
	`

	// Write the file
	err := os.WriteFile(customCfgPath, []byte(cfg), 0o600)
	require.NoError(t, err)

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	err = app.InitConfigFile(tmpDir, customCfgPath)
	require.NoError(t, err)

	// Config file should exist
	require.FileExists(t, path.Join(tmpDir, "config", "config.toml"))

	// read the file
	b, err := os.ReadFile(path.Join(tmpDir, "config", "config.toml"))
	require.NoError(t, err)

	// unmarshal data
	actual := falcon.Config{}
	err = toml.Unmarshal(b, &actual)
	require.NoError(t, err)

	expect := falcon.Config{
		BandChainConfig: band.Config{
			RpcEndpoints: []string{"http://localhost:26657"},
			Timeout:      7,
		},
		TargetChains:           []falcon.TargetChainConfig{},
		CheckingPacketInterval: time.Minute,
	}

	require.Equal(t, expect, actual)
}

func TestInvalidTypeLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// custom config
	customCfgPath := path.Join(tmpDir, "custom.toml")
	// invalid type of value
	invalidType := "'600000'"

	// create new toml config file
	cfg := fmt.Sprintf(`
		target_chains = []
		checking_packet_interval = %v
	
		[bandchain]
		rpc_endpoints = ['http://localhost:26657']
		timeout = 7
	`, invalidType)

	// Write the file
	err := os.WriteFile(customCfgPath, []byte(cfg), 0o600)
	require.NoError(t, err)

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	err = app.InitConfigFile(tmpDir, customCfgPath)

	// should fail when try to unmarshal
	require.ErrorContains(t, err, fmt.Sprintf("LoadConfig file %v error", customCfgPath))

	// file should not be created
	require.NoFileExists(t, path.Join(tmpDir, "config", "config.toml"))
}

func TestInvalidFieldLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// custom config
	customCfgPath := path.Join(tmpDir, "custom.toml")
	// invalid field name
	invalidField := "checking_packet_intervalsss"

	// create new toml config file
	cfg := fmt.Sprintf(`
		target_chains = []
		%v = 60000000000
	
		[bandchain]
		rpc_endpoints = ['http://localhost:26657']
		timeout = 7
	`, invalidField)
	// Write the file
	err := os.WriteFile(customCfgPath, []byte(cfg), 0o600)
	require.NoError(t, err)

	app := falcon.NewApp(nil, nil, tmpDir, false, nil)

	err = app.InitConfigFile(tmpDir, customCfgPath)


	// should fail when try to unmarshal
	require.ErrorContains(t, err, fmt.Sprintf("LoadConfig file %v error", customCfgPath))
	require.ErrorContains(t, err, fmt.Sprintf("invalid field in TOML: %v", invalidField))

	// file should not be created
	require.NoFileExists(t, path.Join(tmpDir, "config", "config.toml"))
}

