package xrpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractKeyName(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected string
		wantErr  bool
	}{
		{
			name:     "standard toml file",
			filePath: "/path/to/key.toml",
			expected: "key",
			wantErr:  false,
		},
		{
			name:     "file no extension",
			filePath: "/path/to/key",
			expected: "key",
			wantErr:  false,
		},
		{
			name:     "relative path",
			filePath: "key.toml",
			expected: "key",
			wantErr:  false,
		},
		{
			name:     "empty filename",
			filePath: "/path/to/.toml",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractKeyName(tt.filePath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestGetXRPLKeyDir(t *testing.T) {
	home := "/home/user"
	chain := "xrpl-mainnet"
	expected := []string{"/home/user", "keys", "xrpl-mainnet", "metadata"}

	got := getXRPLKeyDir(home, chain)
	assert.Equal(t, expected, got)
}
