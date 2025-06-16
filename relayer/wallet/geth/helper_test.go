package geth_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/relayer/wallet/geth"
)

func TestHexToETHAddress(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expect    common.Address
		expectErr bool
	}{
		{
			name:      "valid with 0x prefix",
			input:     "0xe688b84b23f322a994A53dbF8E15FA82CDB71127",
			expect:    common.HexToAddress("0xe688b84b23f322a994A53dbF8E15FA82CDB71127"),
			expectErr: false,
		},
		{
			name:      "valid without 0x",
			input:     "e688b84b23f322a994A53dbF8E15FA82CDB71127",
			expect:    common.HexToAddress("0xe688b84b23f322a994A53dbF8E15FA82CDB71127"),
			expectErr: false,
		},
		{
			name:      "too short",
			input:     "0x1234",
			expectErr: true,
		},
		{
			name:      "invalid character",
			input:     "0xZZZZb84b23f322a994A53dbF8E15FA82CDB71127",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := geth.HexToETHAddress(tc.input)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expect, addr)
			}
		})
	}
}

func TestExtractKeyName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expect    string
		expectErr bool
	}{
		{
			name:      "normal file",
			input:     "/some/dir/foobar.toml",
			expect:    "foobar",
			expectErr: false,
		},
		{
			name:      "no basename",
			input:     "/tmp/.toml",
			expectErr: true,
		},
		{
			name:      "no extension",
			input:     "mykey",
			expect:    "mykey",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := geth.ExtractKeyName(tc.input)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expect, got)
			}
		})
	}
}
