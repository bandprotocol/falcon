package wallet_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

// signalIDFromString right-aligns a string into a [32]byte, matching the
// chain's StringToBytes32 convention (copy to byteArray[32-len(s):]).
func signalIDFromString(s string) [32]byte {
	var id [32]byte
	copy(id[32-len(s):], s)
	return id
}

// The hex payloads below are taken directly from the chain's EncodeTSS output
// (x/tunnel/types/encoding_tss_test.go) and are used as inputs for the
// inverse operation ToTSSPacket.

// TestToTSSPacketFixedPoint verifies that ToTSSPacket correctly decodes a
// FixedPointABI-encoded TSS message produced by the chain.
// Encoded with: sequence=3, prices=[{CS:BAND-USD, 2}], createdAt=123
func TestToTSSPacketFixedPoint(t *testing.T) {
	rawHex := ("cba0ad5a" +
		"0000000000000000000000000000000000000000000000000000000000000020" +
		"0000000000000000000000000000000000000000000000000000000000000003" +
		"0000000000000000000000000000000000000000000000000000000000000060" +
		"000000000000000000000000000000000000000000000000000000000000007b" +
		"0000000000000000000000000000000000000000000000000000000000000001" +
		"00000000000000000000000000000000000000000043533a42414e442d555344" +
		"0000000000000000000000000000000000000000000000000000000000000002")

	msg, err := hex.DecodeString(rawHex)
	require.NoError(t, err)

	payload := wallet.NewTssPayload(msg, nil, nil)
	pkt, err := payload.ToTSSPacket()
	require.NoError(t, err)

	require.Equal(t, uint64(3), pkt.Sequence)
	require.Equal(t, int64(123), pkt.CreatedAt)
	require.Len(t, pkt.RelayPrices, 1)
	require.Equal(t, signalIDFromString("CS:BAND-USD"), pkt.RelayPrices[0].SignalID)
	require.Equal(t, uint64(2), pkt.RelayPrices[0].Price)
}

// TestToTSSPacketTick verifies that ToTSSPacket correctly decodes a
// TickABI-encoded TSS message produced by the chain.
// Encoded with: sequence=3, prices=[{CS:BAND-USD, 2 (→tick 0xf188)}], createdAt=123
func TestToTSSPacketTick(t *testing.T) {
	rawHex := ("db99b2b3" +
		"0000000000000000000000000000000000000000000000000000000000000020" +
		"0000000000000000000000000000000000000000000000000000000000000003" +
		"0000000000000000000000000000000000000000000000000000000000000060" +
		"000000000000000000000000000000000000000000000000000000000000007b" +
		"0000000000000000000000000000000000000000000000000000000000000001" +
		"00000000000000000000000000000000000000000043533a42414e442d555344" +
		"000000000000000000000000000000000000000000000000000000000000f188")

	msg, err := hex.DecodeString(rawHex)
	require.NoError(t, err)

	payload := wallet.NewTssPayload(msg, nil, nil)
	pkt, err := payload.ToTSSPacket()
	require.NoError(t, err)

	require.Equal(t, uint64(3), pkt.Sequence)
	require.Equal(t, int64(123), pkt.CreatedAt)
	require.Len(t, pkt.RelayPrices, 1)
	require.Equal(t, signalIDFromString("CS:BAND-USD"), pkt.RelayPrices[0].SignalID)
	require.Equal(t, uint64(0xf188), pkt.RelayPrices[0].Price)
}

// TestToTSSPacketTooShort verifies that a message shorter than the 4-byte
// prefix returns an error.
func TestToTSSPacketTooShort(t *testing.T) {
	for _, msg := range [][]byte{
		{},
		{0x01},
		{0x01, 0x02, 0x03},
	} {
		payload := wallet.NewTssPayload(msg, nil, nil)
		_, err := payload.ToTSSPacket()
		require.Error(t, err)
		require.ErrorContains(t, err, "tss message should have at least")
	}
}

// TestToTSSPacketInvalidABI verifies that garbage bytes after the 4-byte
// prefix produce an unpack error.
func TestToTSSPacketInvalidABI(t *testing.T) {
	msg := []byte{0x00, 0x00, 0x00, 0x00, 0xDE, 0xAD, 0xBE, 0xEF}
	payload := wallet.NewTssPayload(msg, nil, nil)
	_, err := payload.ToTSSPacket()
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to unpack ABI data")
}

// TestBytes32ToString verifies the helper used to recover SignalID strings.
func TestBytes32ToString(t *testing.T) {
	tests := []struct {
		name     string
		input    [32]byte
		expected string
	}{
		{
			name:     "right-aligned short string",
			input:    signalIDFromString("CS:BAND-USD"),
			expected: "CS:BAND-USD",
		},
		{
			name:     "single char",
			input:    signalIDFromString("X"),
			expected: "X",
		},
		{
			name:     "full 32-byte string",
			input:    signalIDFromString("12345678901234567890123456789012"),
			expected: "12345678901234567890123456789012",
		},
		{
			name:     "all zero bytes returns empty string",
			input:    [32]byte{},
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := wallet.Bytes32ToString(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}
}
