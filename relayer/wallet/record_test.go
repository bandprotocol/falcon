package wallet_test

import (
	stdos "os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/internal/os"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

func TestLoadKeyRecords(t *testing.T) {
	type file struct {
		name string
		data []byte
	}
	tests := []struct {
		name            string
		setupFiles      []file
		lookupPath      string
		expectedRecords map[string]wallet.KeyRecord
		expectErr       bool
	}{
		{
			name: "success two records",
			setupFiles: []file{
				{"alice.toml", []byte(`type = "local"` + "\n" + `address = "0xAAA"`)},
				{"bob.toml", []byte(`type = "remote"` + "\n" + `address = "0xBBB"` + "\n" + `url = "http://example.com"`)},
			},
			expectedRecords: map[string]wallet.KeyRecord{
				"alice": {Type: "local", Address: "0xAAA"},
				"bob":   {Type: "remote", Address: "0xBBB", Url: "http://example.com"},
			},
		},
		{
			name:            "empty (nonexistent dir)",
			lookupPath:      path.Join("does_not_exist_dir"),
			expectedRecords: map[string]wallet.KeyRecord{},
		},
		{
			name: "invalid toml",
			setupFiles: []file{
				{"bad.toml", []byte("not = valid = toml")},
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			for _, f := range tc.setupFiles {
				err := os.Write(
					f.data,
					[]string{tmpDir, f.name},
				)
				require.NoError(t, err)
			}

			lookup := tmpDir
			if tc.lookupPath != "" {
				lookup = path.Join(tmpDir, tc.lookupPath)
			}

			records, err := wallet.LoadKeyRecords(lookup)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedRecords, records)
		})
	}
}

func TestNewKeyRecord(t *testing.T) {
	key := "test-key"
	record := wallet.NewKeyRecord("local", "address", "url", &key)

	assert.Equal(t, "local", record.Type)
	assert.Equal(t, "address", record.Address)
	assert.Equal(t, "url", record.Url)
	assert.Equal(t, &key, record.Key)
}

func TestLoadKeyRecord(t *testing.T) {
	tmpDir, err := stdos.MkdirTemp("", "wallet-test")
	require.NoError(t, err)
	defer stdos.RemoveAll(tmpDir)

	keyName := "test-signer"
	content := `
type = "local"
address = "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2"
`
	err = stdos.WriteFile(filepath.Join(tmpDir, keyName+".toml"), []byte(content), 0o600)
	require.NoError(t, err)

	records, err := wallet.LoadKeyRecords(tmpDir)
	require.NoError(t, err)
	assert.Len(t, records, 1)
	assert.Contains(t, records, keyName)

	record := records[keyName]
	assert.Equal(t, "local", record.Type)
	assert.Equal(t, "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2", record.Address)
}
