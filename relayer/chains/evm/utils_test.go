package evm_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/relayer/chains/evm"
)

func TestHexToAddress(t *testing.T) {
	type TestCases struct {
		input  string
		output string
		err    error
	}

	testcases := []TestCases{
		{
			input:  "0x1234567890123456789012345678901234567890",
			output: "0x1234567890123456789012345678901234567890",
			err:    nil,
		},
		{
			input:  "1234567890123456789012345678901234567890",
			output: "0x1234567890123456789012345678901234567890",
			err:    nil,
		},
		{
			input:  "0xe688b84b23f322a994A53dbF8E15FA82CDB71127",
			output: "0xe688b84b23f322a994A53dbF8E15FA82CDB71127",
			err:    nil,
		},
		{
			input:  "e688b84b23f322a994A53dbF8E15FA82CDB71127",
			output: "0xe688b84b23f322a994A53dbF8E15FA82CDB71127",
			err:    nil,
		},
		{
			input:  "0xE688B84B23F322A994A53DBF8E15FA82CDB71127",
			output: "0xe688b84b23f322a994A53dbF8E15FA82CDB71127",
			err:    nil,
		},
		{
			input:  "0x123",
			output: "",
			err:    fmt.Errorf("invalid address: 0x123"),
		},
		{
			input:  "0x123456789012345678901234567890123456789Z",
			output: "",
			err:    fmt.Errorf("invalid address: 0x123456789012345678901234567890123456789Z"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			addr, err := evm.HexToAddress(tc.input)
			if tc.err != nil {
				require.ErrorContains(t, tc.err, err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.output, addr.Hex())
			}
		})
	}
}

func TestStripPrivateKeyPrefix(t *testing.T) {
	tests := []struct {
		name       string
		privateKey string
		expected   string
	}{
		{
			name:       "With 0x prefix",
			privateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			expected:   "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		{
			name:       "Without 0x prefix",
			privateKey: "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			expected:   "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evm.StripPrivateKeyPrefix(tt.privateKey)
			require.Equal(t, tt.expected, result, "unexpected result for privateKey %s", tt.privateKey)
		})
	}
}

func TestMultiplyBigIntWithFloat64(t *testing.T) {
	tests := []struct {
		name       string
		input      *big.Int
		multiplier float64
		expected   *big.Int
	}{
		{
			name:       "Multiply positive value",
			input:      big.NewInt(100),
			multiplier: 2.5,
			expected:   big.NewInt(250),
		},
		{
			name:       "Multiply by zero",
			input:      big.NewInt(100),
			multiplier: 0,
			expected:   big.NewInt(0),
		},
		{
			name:       "Multiply negative value",
			input:      big.NewInt(-100),
			multiplier: 1.5,
			expected:   big.NewInt(-150),
		},
		{
			name:       "Multiply large number",
			input:      big.NewInt(1e6),
			multiplier: 1.1,
			expected:   big.NewInt(1100000),
		},
		{
			name:       "Multiply by fractional multiplier",
			input:      big.NewInt(100),
			multiplier: 0.333,
			expected:   big.NewInt(33), // Rounded down
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evm.MultiplyBigIntWithFloat64(tt.input, tt.multiplier)
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("MultiplyBigIntWithFloat64(%v, %f) = %v, want %v",
					tt.input, tt.multiplier, result, tt.expected)
			}
		})
	}
}
