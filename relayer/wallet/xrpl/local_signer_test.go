package xrpl_test

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/xrpl"
)

// validTssMessage returns a FixedPointABI-encoded TSS message that contains
// one relay price for "CS:BAND-USD". The hex is taken from the chain's own
// encoding_tss_test.go to guarantee correctness:
//
//	sequence=3, prices=[{CS:BAND-USD, price=2}], createdAt=123
func validTssMessage(t *testing.T) []byte {
	t.Helper()
	rawHex := (
	// 52 bytes of zero padding before the 4-byte selector
	// (EncoderABIPrefixLength = 52 bytes + 4-byte selector = 56 bytes total)
	"0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000" +
		"cba0ad5a" +
		// ABI-encoded payload
		"0000000000000000000000000000000000000000000000000000000000000020" +
		"0000000000000000000000000000000000000000000000000000000000000003" +
		"0000000000000000000000000000000000000000000000000000000000000060" +
		"000000000000000000000000000000000000000000000000000000000000007b" +
		"0000000000000000000000000000000000000000000000000000000000000001" +
		"00000000000000000000000000000000000000000043533a42414e442d555344" +
		"0000000000000000000000000000000000000000000000000000000000000002")
	b, err := hex.DecodeString(rawHex)
	require.NoError(t, err)
	return b
}

func TestLocalSigner(t *testing.T) {
	// Root account secret for testing
	seed := "sEdVeuhfwHB6dMxgSBccJ7ZYGyLfySa"
	w, err := xrplwallet.FromSecret(seed)
	require.NoError(t, err)

	name := "test-local-signer"
	signer := xrpl.NewLocalSigner(name, &w)

	assert.Equal(t, name, signer.GetName())
	assert.Equal(t, w.ClassicAddress.String(), signer.GetAddress())

	privKey, err := signer.ExportPrivateKey()
	assert.NoError(t, err)
	assert.Equal(t, w.PrivateKey, privKey)

	// Test Sign
	signerPayload := xrpl.NewSignerPayload(w.ClassicAddress.String(), 0, "100", 1)

	payloadBytes, err := json.Marshal(signerPayload)
	require.NoError(t, err)

	tssPayload := wallet.NewTssPayload(validTssMessage(t), nil, nil)
	signedBlob, err := signer.Sign(payloadBytes, tssPayload)
	assert.NoError(t, err)
	assert.NotEmpty(t, signedBlob)
}
