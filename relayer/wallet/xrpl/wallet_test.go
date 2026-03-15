package xrpl_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/falcon/relayer/wallet/xrpl"
)

func TestXRPLWallet(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "xrpl-wallet-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpHome)

	passphrase := "test-passphrase"
	chainName := "xrpl-test"

	// Create necessary directories
	metadataDir := filepath.Join(tmpHome, "keys", chainName, "metadata")
	err = os.MkdirAll(metadataDir, 0o755)
	require.NoError(t, err)

	// Step 1: Initialize empty wallet
	w, err := xrpl.NewWallet(passphrase, tmpHome, chainName)
	require.NoError(t, err)
	assert.Empty(t, w.GetSigners())

	// Step 2: Save key by mnemonic
	name2 := "mnemonic-key"
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	addr2, err := w.SaveByMnemonic(name2, mnemonic, 144, 0, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, addr2)
	assert.True(t, w.IsAddressExist(addr2))

	// Step 3: Save remote signer key
	name3 := "remote-key"
	addr3 := "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2"
	url := "localhost:50051"
	err = w.SaveRemoteSignerKey(name3, addr3, url, "")
	require.NoError(t, err)
	assert.True(t, w.IsAddressExist(addr3))

	// Step 4: Verify retrieval
	signers := w.GetSigners()
	require.Len(t, signers, 2)

	signer, ok := w.GetSigner(name2)
	require.True(t, ok)
	assert.Equal(t, addr2, signer.GetAddress())

	// Step 5: Delete key
	err = w.DeleteKey(name2)
	assert.NoError(t, err)
	assert.False(t, w.IsAddressExist(addr2))
	assert.Len(t, w.GetSigners(), 1)

	// Step 6: Test re-initialization (should load saved keys)
	// We need to re-create the wallet object to simulate restart
	w2, err := xrpl.NewWallet(passphrase, tmpHome, chainName)
	assert.NoError(t, err)
	assert.Len(t, w2.GetSigners(), 1)
	assert.True(t, w2.IsAddressExist(addr3))
}
