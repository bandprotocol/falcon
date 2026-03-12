package geth_test

import (
	"path"
	"testing"

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
