package secret_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/relayer/wallet/secret"
)

func TestNewSignerPayload(t *testing.T) {
	payload := secret.NewSignerPayload(
		"secret1sender",
		"secret1contract",
		"testing",
		1,
		2,
		3000000,
		"0.1uscrt",
		"test memo",
		"1234abcd",
		"pubkey123",
	)

	assert.Equal(t, "secret1sender", payload.Sender)
	assert.Equal(t, "secret1contract", payload.ContractAddress)
	assert.Equal(t, "testing", payload.ChainID)
	assert.Equal(t, uint64(1), payload.AccountNumber)
	assert.Equal(t, uint64(2), payload.Sequence)
	assert.Equal(t, uint64(3000000), payload.GasLimit)
	assert.Equal(t, "0.1uscrt", payload.GasPrices)
	assert.Equal(t, "test memo", payload.Memo)
	assert.Equal(t, "1234abcd", payload.CodeHash)
	assert.Equal(t, "pubkey123", payload.ChainPubkey)
}

func TestSecretWalletRemoteOnly(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "secret-wallet-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpHome)

	chainName := "secret-test"
	passphrase := "test-passphrase"

	w, err := secret.NewWallet(passphrase, tmpHome, chainName)
	require.NoError(t, err)
	assert.Empty(t, w.GetSigners())

	// Local key import should not be supported.
	_, err = w.SaveByMnemonic("mnemonic-1", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about", 1, 0, 0)
	assert.Error(t, err)

	// Remote signer keys are supported.
	remoteName := "remote-key"
	remoteAddr := "secret1x1234567890abcdef"
	remoteURL := "localhost:50051"
	assert.NoError(t, w.SaveRemoteSignerKey(remoteName, remoteAddr, remoteURL, "api-key"))
	assert.True(t, w.IsAddressExist(remoteAddr))
	assert.Len(t, w.GetSigners(), 1)

	signer, ok := w.GetSigner(remoteName)
	require.True(t, ok)
	assert.Equal(t, remoteAddr, signer.GetAddress())
}
