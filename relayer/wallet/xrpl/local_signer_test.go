package xrpl_test

import (
	"encoding/json"
	"testing"

	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/xrpl"
)

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
	signerPayload := xrpl.SignerPayload{
		Account: w.ClassicAddress.String(),
		Fee:     "100",
	}

	payloadBytes, err := json.Marshal(signerPayload)
	require.NoError(t, err)

	signedBlob, err := signer.Sign(payloadBytes, wallet.TssPayload{})
	assert.NoError(t, err)
	assert.NotEmpty(t, signedBlob)
}
