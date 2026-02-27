package chains

import (
	"github.com/bandprotocol/falcon/relayer/wallet"
)

// LoadSigners returns the Signer channel with all configured wallet signers.
func LoadSigners(w wallet.Wallet) chan wallet.Signer {
	signers := w.GetSigners()
	signerChannel := make(chan wallet.Signer, len(signers))

	for _, signer := range signers {
		signerChannel <- signer
	}

	return signerChannel
}
