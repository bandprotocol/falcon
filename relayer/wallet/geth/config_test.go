package geth_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/internal/os"
	"github.com/bandprotocol/falcon/relayer/wallet/geth"
)

func TestUnmarshalSignerRecord(t *testing.T) {
	type file struct {
		name string
		data []byte
	}
	tests := []struct {
		name            string
		setupFiles      []file
		lookupPath      string
		expectedRecords map[string]geth.SignerRecord
		expectErr       bool
	}{
		{
			name: "success two records",
			setupFiles: []file{
				{"alice.toml", []byte(`address="0xAAA"` + "\n" + `type="local"`)},
				{"bob.toml", []byte(`address="0xBBB"` + "\n" + `type="remote"`)},
			},
			expectedRecords: map[string]geth.SignerRecord{
				"alice": {Address: "0xAAA", Type: "local"},
				"bob":   {Address: "0xBBB", Type: "remote"},
			},
		},
		{
			name:            "empty (nonexistent dir)",
			lookupPath:      path.Join("does_not_exist_dir"),
			expectedRecords: map[string]geth.SignerRecord{},
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

			records, err := geth.LoadSignerRecord(lookup)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedRecords, records)
		})
	}
}

func TestUnmarshalRemoteSignerRecord(t *testing.T) {
	type file struct {
		name string
		data []byte
	}
	tests := []struct {
		name            string
		setupFiles      []file
		lookupPath      string
		expectedRecords map[string]geth.RemoteSignerRecord
		expectErr       bool
	}{
		{
			name: "success two remote",
			setupFiles: []file{
				{"node1.toml", []byte(`url="http://example.com:1234"`)},
				{"node2.toml", []byte(`url="https://node.local"`)},
			},
			expectedRecords: map[string]geth.RemoteSignerRecord{
				"node1": {Url: "http://example.com:1234"},
				"node2": {Url: "https://node.local"},
			},
		},
		{
			name:            "empty (nonexistent dir)",
			lookupPath:      path.Join("no_such"),
			expectedRecords: map[string]geth.RemoteSignerRecord{},
		},
		{
			name: "invalid toml",
			setupFiles: []file{
				{"bad.toml", []byte("url :=://")},
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

			records, err := geth.LoadRemoteSignerRecord(lookup)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedRecords, records)
		})
	}
}
