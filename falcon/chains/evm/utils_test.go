package evm_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/falcon/chains/evm"
)

func TestHexToAddress(t *testing.T) {
	type TestCases struct {
		input  string
		output string
		err    error
	}

	testcases := []TestCases{
		{"0x1234567890123456789012345678901234567890", "0x1234567890123456789012345678901234567890", nil},
		{"1234567890123456789012345678901234567890", "0x1234567890123456789012345678901234567890", nil},
		{"0xe688b84b23f322a994A53dbF8E15FA82CDB71127", "0xe688b84b23f322a994A53dbF8E15FA82CDB71127", nil},
		{"e688b84b23f322a994A53dbF8E15FA82CDB71127", "0xe688b84b23f322a994A53dbF8E15FA82CDB71127", nil},
		{"0xE688B84B23F322A994A53DBF8E15FA82CDB71127", "0xe688b84b23f322a994A53dbF8E15FA82CDB71127", nil},
		{"0x123", "", fmt.Errorf("invalid address: 0x123")},
		{
			"0x123456789012345678901234567890123456789Z",
			"",
			fmt.Errorf("invalid address: 0x123456789012345678901234567890123456789Z"),
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
