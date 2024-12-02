package evm_test

import (
	"fmt"
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

func TestConvertPrivateKeyStrToHexWithPrefixWith0xPrefix(t *testing.T) {
	privateKey := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	expected := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	require.Equal(t, expected, evm.ConvertPrivateKeyStrToHex(privateKey))
}

func TestConvertPrivateKeyStrToHexWithPrefixWithout0xPrefix(t *testing.T) {
	privateKey := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	require.Equal(t, privateKey, evm.ConvertPrivateKeyStrToHex(privateKey))
}
