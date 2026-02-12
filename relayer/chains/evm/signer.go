package evm

import (
	"github.com/bandprotocol/falcon/relayer/chains"
)

// LoadSigners initializes the Signer channel with all configured wallet signers.
func (cp *EVMChainProvider) LoadSigners() error {
	cp.FreeSigners = chains.LoadSigners(cp.Wallet)
	return nil
}
