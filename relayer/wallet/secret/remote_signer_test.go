package secret_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/secret"
)

func TestSecretRemoteSignerSign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFkmsClient := mocks.NewMockFkmsServiceClient(ctrl)
	name := "test-remote"
	address := "secret1x1234567890abcdef"
	url := "localhost:50051"
	apiKey := "test-api-key"

	signerIface, err := secret.NewRemoteSigner(name, address, url, apiKey)
	require.NoError(t, err)

	signer, ok := signerIface.(*secret.RemoteSigner)
	require.True(t, ok)

	signer.FkmsClient = mockFkmsClient

	assert.Equal(t, name, signer.GetName())
	assert.Equal(t, address, signer.GetAddress())

	_, err = signer.ExportPrivateKey()
	assert.Error(t, err)

	signerPayload := secret.SignerPayload{
		Sender:          "secret1sender",
		ContractAddress: "secret1contract",
		ChainID:         "testing",
		AccountNumber:   1,
		Sequence:        2,
		GasLimit:        3000000,
		GasPrices:       "0.1uscrt",
		Memo:            "test memo",
		CodeHash:        "1234abcd",
		ChainPubkey:     "pubkey123",
	}

	payload, err := json.Marshal(signerPayload)
	require.NoError(t, err)

	tssPayload := wallet.TssPayload{
		TssMessage: []byte("tss-msg"),
		RandomAddr: []byte("random-addr"),
		Signature:  []byte("signature"),
	}

	expectedTxBlob := []byte("signed-secret-tx")

	mockFkmsClient.EXPECT().SignSecret(
		gomock.Any(),
		&fkmsv1.SignSecretRequest{
			SignerPayload: &fkmsv1.SecretSignerPayload{
				Sender:          signerPayload.Sender,
				ContractAddress: signerPayload.ContractAddress,
				ChainId:         signerPayload.ChainID,
				AccountNumber:   signerPayload.AccountNumber,
				Sequence:        signerPayload.Sequence,
				GasLimit:        signerPayload.GasLimit,
				GasPrices:       signerPayload.GasPrices,
				Memo:            signerPayload.Memo,
				CodeHash:        signerPayload.CodeHash,
				ChainPubkey:     signerPayload.ChainPubkey,
			},
			Tss: &fkmsv1.Tss{
				Message:    tssPayload.TssMessage,
				RandomAddr: tssPayload.RandomAddr,
				SignatureS: tssPayload.Signature,
			},
		},
	).Return(&fkmsv1.SignSecretResponse{TxBlob: expectedTxBlob}, nil)

	signedBlob, err := signer.Sign(payload, tssPayload)
	require.NoError(t, err)
	assert.Equal(t, expectedTxBlob, signedBlob)
}
