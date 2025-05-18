package evm

import (
	"github.com/bandprotocol/falcon/relayer/wallet"
)

// LoadSigners initializes the Signer channel with all configured wallet signers.
func (cp *EVMChainProvider) LoadSigners() error {
	signers := cp.Wallet.GetSigners()
	signerChannels := make(chan wallet.Signer, len(signers))

	for _, signer := range signers {
		signerChannels <- signer
	}

	cp.Signer = signerChannels
	return nil
}
