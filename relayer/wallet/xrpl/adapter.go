package xrpl

import (
	"fmt"

	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

const xrplDefaultCoinType = 144

// Adapter implements wallet.WalletAdapter for XRPL chains.
type Adapter struct{}

var _ wallet.WalletAdapter = (*Adapter)(nil)

// LoadSigner reconstructs a Signer from a persisted KeyRecord and the secret
// retrieved from the shared keyring.
func (a *Adapter) LoadSigner(name string, record wallet.KeyRecord, secret string) (wallet.Signer, error) {
	switch record.Type {
	case wallet.MnemonicSignerType:
		ms, err := wallet.DecodeMnemonicSecret(secret)
		if err != nil {
			return nil, err
		}
		if ms.CoinType != xrplDefaultCoinType || ms.Account != 0 || ms.Index != 0 {
			return nil, fmt.Errorf("XRPL only supports derivation path m/44'/144'/0'/0/0")
		}
		mnWallet, err := xrplwallet.FromMnemonic(ms.Mnemonic)
		if err != nil {
			return nil, err
		}
		return NewLocalSigner(name, mnWallet), nil
	case wallet.PrivKeySignerType:
		return nil, fmt.Errorf("XRPL does not support private key import")
	case wallet.RemoteSignerType:
		return NewRemoteSigner(name, record.Address, record.URL, secret)
	default:
		return nil, fmt.Errorf("unsupported signer type: %s for key %s", record.Type, name)
	}
}
