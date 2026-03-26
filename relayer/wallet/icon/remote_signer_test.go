package icon_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bandprotocol/falcon/internal/relayertest/mocks"
	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/icon"
)

func TestIconRemoteSignerSign(t *testing.T) {
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
		Relayer:         "blue-relayer",
		ContractAddress: "cxdeadbeef",
		StepLimit:       12345,
		NetworkID:       "testnet",
	}

	payload, err := json.Marshal(signerPayload)
	require.NoError(t, err)

	tssPayload := wallet.TssPayload{
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
				Message:    tssPayload.TssMessage,
				RandomAddr: tssPayload.RandomAddr,
				SignatureS: tssPayload.Signature,
			},
		},
	).Return(&fkmsv1.SignIconResponse{TxParams: expectedTxParams}, nil)

	signedBlob, err := signer.Sign(payload, tssPayload)
	require.NoError(t, err)
	assert.Equal(t, expectedTxParams, signedBlob)
}
