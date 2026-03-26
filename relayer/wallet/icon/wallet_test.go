package icon_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/icon"
)

func TestNewSignerPayload(t *testing.T) {
	payload := icon.NewSignerPayload("my-relayer", "cx1234567890", 100000, "mainnet")

	assert.Equal(t, "my-relayer", payload.Relayer)
	assert.Equal(t, "cx1234567890", payload.ContractAddress)
	assert.Equal(t, uint64(100000), payload.StepLimit)
	assert.Equal(t, "mainnet", payload.NetworkID)
}

func TestIconRemoteSigner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFkmsClient := mocks.NewMockFkmsServiceClient(ctrl)

	name := "test-remote"
	address := "cx1234567890abcdef"
	url := "localhost:50051"
	apiKey := "test-api-key"

	signerIface, err := icon.NewRemoteSigner(name, address, url, apiKey)
	require.NoError(t, err)

	signer, ok := signerIface.(*icon.RemoteSigner)
	require.True(t, ok)

	signer.FkmsClient = mockFkmsClient

	assert.Equal(t, name, signer.GetName())
	assert.Equal(t, address, signer.GetAddress())

	_, err = signer.ExportPrivateKey()
	assert.Error(t, err)

	signerPayload := icon.SignerPayload{
		Relayer:         "my-relayer",
		ContractAddress: "cx1234567890abcdef",
		StepLimit:       5000,
		NetworkID:       "mainnet",
	}

	payload, err := json.Marshal(signerPayload)
	require.NoError(t, err)

	preSignPayload := &wallet.TssPayload{
		TssMessage: []byte("tss-msg"),
		RandomAddr: []byte("random-addr"),
		Signature:  []byte("signature"),
	}

	expectedTxParams := []byte("signed-icon-tx")

	mockFkmsClient.EXPECT().SignIcon(
		gomock.Any(),
		&fkmsv1.SignIconRequest{
			SignerPayload: &fkmsv1.IconSignerPayload{
				Relayer:         signerPayload.Relayer,
				ContractAddress: signerPayload.ContractAddress,
				StepLimit:       signerPayload.StepLimit,
				NetworkId:       signerPayload.NetworkID,
			},
			Tss: &fkmsv1.Tss{
				Message:    preSignPayload.TssMessage,
				RandomAddr: preSignPayload.RandomAddr,
				SignatureS: preSignPayload.Signature,
			},
		},
	).Return(&fkmsv1.SignIconResponse{TxParams: expectedTxParams}, nil)

	signedBlob, err := signer.Sign(payload, *preSignPayload)
	require.NoError(t, err)
	assert.Equal(t, expectedTxParams, signedBlob)
}

func TestIconWalletRemoteOnly(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "icon-wallet-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpHome)

	chainName := "icon-test"
	passphrase := "test-passphrase"

	w, err := icon.NewWallet(passphrase, tmpHome, chainName)
	require.NoError(t, err)
	assert.Empty(t, w.GetSigners())

	// Local key import should not be supported.
	_, err = w.SaveByMnemonic("mnemonic-1", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about", 1, 0, 0)
	assert.Error(t, err)

	// Remote signer keys are supported.
	remoteName := "remote-key"
	remoteAddr := "cx1234567890abcdef"
	remoteURL := "localhost:50051"
	assert.NoError(t, w.SaveRemoteSignerKey(remoteName, remoteAddr, remoteURL, "api-key"))
	assert.True(t, w.IsAddressExist(remoteAddr))
	assert.Len(t, w.GetSigners(), 1)

	signer, ok := w.GetSigner(remoteName)
	require.True(t, ok)
	assert.Equal(t, remoteAddr, signer.GetAddress())
}
