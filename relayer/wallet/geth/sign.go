package geth

import (
	"fmt"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

func SignEvmTx(signer wallet.Signer, data []byte) ([]byte, error) {
	switch s := signer.(type) {
	case *LocalSigner:
		return s.localSign(data)
	case *RemoteSigner:
		return s.remoteSign(data)
	default:
		return []byte{}, fmt.Errorf("unsupported signer type: %T", signer)
	}
}
