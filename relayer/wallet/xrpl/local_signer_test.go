package xrpl_test

import (
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
	// Valid XRPL transaction hex for binarycodec.Decode
	txHex := "120000240000000A201B00000000614000000000000064684000000000000000732103AD23396EB905659472C15EE780A9001D48A5E784EF6E09E96BE0A43AF24647318114AD23396EB905659472C15EE780A9001D48A5E784"

	signedBlob, err := signer.Sign([]byte(txHex), wallet.PreSignPayload{})
	assert.NoError(t, err)
	assert.NotEmpty(t, signedBlob)
}
