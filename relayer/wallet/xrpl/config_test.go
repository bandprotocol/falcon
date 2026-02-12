package xrpl_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/relayer/wallet/xrpl"
)

func TestNewKeyRecord(t *testing.T) {
	key := "test-key"
	record := xrpl.NewKeyRecord("local", "address", "url", &key, "seed")

	assert.Equal(t, "local", record.Type)
	assert.Equal(t, "address", record.Address)
	assert.Equal(t, "url", record.Url)
	assert.Equal(t, &key, record.Key)
	assert.Equal(t, "seed", record.SaveMethod)
}

func TestLoadKeyRecord(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "xrpl-wallet-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a dummy key record file
	keyName := "test-signer"
	content := `
type = "local"
address = "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2"
save_method = "seed"
`
	err = os.WriteFile(filepath.Join(tmpDir, keyName+".toml"), []byte(content), 0o600)
	require.NoError(t, err)

	records, err := xrpl.LoadKeyRecord(tmpDir)
	require.NoError(t, err)
	assert.Len(t, records, 1)
	assert.Contains(t, records, keyName)

	record := records[keyName]
	assert.Equal(t, "local", record.Type)
	assert.Equal(t, "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2", record.Address)
	assert.Equal(t, "seed", record.SaveMethod)
}
