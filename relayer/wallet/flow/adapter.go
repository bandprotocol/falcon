package flow

import (
	"fmt"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

// Adapter implements wallet.WalletAdapter for Flow chains.
// Flow is remote-only — local key import is not supported.
type Adapter struct{}

var _ wallet.WalletAdapter = (*Adapter)(nil)

// LoadSigner reconstructs a Signer from a persisted KeyRecord and the secret
// retrieved from the shared keyring. Local signer types (private key and mnemonic)
// are not supported for Flow; only remote signing via KMS is allowed.
func (a *Adapter) LoadSigner(name string, record wallet.KeyRecord, secret string) (wallet.Signer, error) {
	switch record.Type {
	case wallet.RemoteSignerType:
		return NewRemoteSigner(name, record.Address, record.URL, secret)
	case wallet.PrivKeySignerType:
		return nil, fmt.Errorf("flow does not support private key import")
	case wallet.MnemonicSignerType:
		return nil, fmt.Errorf("flow does not support mnemonic import")
	default:
		return nil, fmt.Errorf("unsupported signer type: %s for key %s", record.Type, name)
	}
}
