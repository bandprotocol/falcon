package evm

import (
	"github.com/bandprotocol/falcon/relayer/wallet"
)

// LoadSigners initializes the Signer channel with all configured wallet signers.
func (cp *EVMChainProvider) LoadSigners() error {
	signers := cp.Wallet.GetSigners()
	signerChannel := make(chan wallet.Signer, len(signers))

	for _, signer := range signers {
		signerChannel <- signer
	}

	cp.FreeSigners = signerChannel
	return nil
}
