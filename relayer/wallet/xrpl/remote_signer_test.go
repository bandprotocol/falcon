package xrpl_test

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/xrpl"
)

func TestRemoteSigner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFkmsClient := mocks.NewMockFkmsServiceClient(ctrl)

	name := "test-remote-signer"
	address := "rHb9CJAW8f5rjR5juUs6K3mJtr47MS9f2"
	url := "localhost:50051"
	apiKey := "test-api-key"

	signer, err := xrpl.NewRemoteSigner(name, address, url, apiKey)
	assert.NoError(t, err)

	// Overwrite the client with our mock
	signer.FkmsClient = mockFkmsClient

	assert.Equal(t, name, signer.GetName())
	assert.Equal(t, address, signer.GetAddress())

	_, err = signer.ExportPrivateKey()
	assert.Error(t, err)

	// Test Sign
	signerPayload := xrpl.SignerPayload{
		Account: address,
		Fee:     "100",
	}
	preSignPayload := &wallet.TssPayload{
		TssMessage: []byte("tss-msg"),
		RandomAddr: []byte("random-addr"),
		Signature:  []byte("signature"),
	}

	expectedTxBlob := []byte("signed-tx-blob")
	mockFkmsClient.EXPECT().SignXrpl(
		gomock.Any(),
		&fkmsv1.SignXrplRequest{
			SignerPayload: &fkmsv1.XrplSignerPayload{
				Account: signerPayload.Account,
				Fee:     signerPayload.Fee,
			},
			Tss: &fkmsv1.Tss{
				Message:    preSignPayload.TssMessage,
				RandomAddr: preSignPayload.RandomAddr,
				SignatureS: preSignPayload.Signature,
			},
		},
	).Return(&fkmsv1.SignXrplResponse{TxBlob: expectedTxBlob}, nil)

	payload, err := json.Marshal(signerPayload)
	assert.NoError(t, err)

	signedBlob, err := signer.Sign(payload, *preSignPayload)
	assert.NoError(t, err)
	assert.Equal(t, hex.EncodeToString(expectedTxBlob), string(signedBlob))
}
