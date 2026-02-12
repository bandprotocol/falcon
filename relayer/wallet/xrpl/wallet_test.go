package xrpl_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bandprotocol/falcon/relayer/wallet/xrpl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXRPLWallet(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "xrpl-wallet-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpHome)

	passphrase := "test-passphrase"
	chainName := "xrpl-test"

	// Create necessary directories
	metadataDir := filepath.Join(tmpHome, "keys", chainName, "metadata")
	err = os.MkdirAll(metadataDir, 0755)
	require.NoError(t, err)

	privDir := filepath.Join(tmpHome, "keys", chainName, "priv")
	err = os.MkdirAll(privDir, 0755)
	require.NoError(t, err)

	// Step 1: Initialize empty wallet
	w, err := xrpl.NewXRPLWallet(passphrase, tmpHome, chainName)
	require.NoError(t, err)
	assert.Empty(t, w.GetSigners())

	// Step 2: Save key by family seed
	name1 := "seed-key"
	seed := "sEdVeuhfwHB6dMxgSBccJ7ZYGyLfySa"
	addr1, err := w.SaveByFamilySeed(name1, seed)
	require.NoError(t, err)
	assert.NotEmpty(t, addr1)
	assert.True(t, w.IsAddressExist(addr1))

	// Step 3: Save key by mnemonic
	name2 := "mnemonic-key"
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	addr2, err := w.SaveByMnemonic(name2, mnemonic, 144, 0, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, addr2)
	assert.True(t, w.IsAddressExist(addr2))

	// Step 4: Save remote signer key
	name3 := "remote-key"
	addr3 := "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2"
	url := "localhost:50051"
	err = w.SaveRemoteSignerKey(name3, addr3, url, nil)
	require.NoError(t, err)
	assert.True(t, w.IsAddressExist(addr3))

	// Step 5: Verify retrieval
	signers := w.GetSigners()
	require.Len(t, signers, 3)

	signer, ok := w.GetSigner(name1)
	require.True(t, ok)
	assert.Equal(t, addr1, signer.GetAddress())

	// Step 6: Delete key
	err = w.DeleteKey(name1)
	assert.NoError(t, err)
	assert.False(t, w.IsAddressExist(addr1))
	assert.Len(t, w.GetSigners(), 2)

	// Step 7: Test re-initialization (should load saved keys)
	// We need to re-create the wallet object to simulate restart
	w2, err := xrpl.NewXRPLWallet(passphrase, tmpHome, chainName)
	assert.NoError(t, err)
	assert.Len(t, w2.GetSigners(), 2)
	assert.True(t, w2.IsAddressExist(addr2))
	assert.True(t, w2.IsAddressExist(addr3))
}
