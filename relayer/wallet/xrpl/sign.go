package xrpl

import (
	"fmt"

	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

type SignerPayload struct {
	SignalPrices    []bandtypes.SignalPrice
	Account         string
	OracleId        uint64
	Fee             string
	Sequence        uint64
	LastUpdatedTime uint64
}

func NewSignerPayload(signalPrices []bandtypes.SignalPrice, account string, oracleId uint64, fee string, sequence uint64, lastUpdatedTime uint64) SignerPayload {
	return SignerPayload{
		SignalPrices:    signalPrices,
		Account:         account,
		OracleId:        oracleId,
		Fee:             fee,
		Sequence:        sequence,
		LastUpdatedTime: lastUpdatedTime,
	}
}

func SignXrplTx(signer wallet.Signer, signerPayload SignerPayload, tssPayload wallet.TssPayload) (string, error) {
	switch s := signer.(type) {
	case *LocalSigner:
		return s.localSign(signerPayload)
	case *RemoteSigner:
		return s.remoteSign(signerPayload, tssPayload)
	default:
		return "", fmt.Errorf("unsupported signer type: %T", signer)
	}
}
