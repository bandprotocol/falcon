package xrpl

import (
	"fmt"

	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

const (
	xrplDefaultCoinType = 144
)

// Adapter implements wallet.WalletAdapter for XRPL chains.
type Adapter struct {
	passphrase string
	homePath   string
	chainName  string
}

var _ wallet.WalletAdapter = (*Adapter)(nil)

// DeriveFromPrivateKey is not supported for XRPL.
func (a *Adapter) DeriveFromPrivateKey(name, privateKey string) (wallet.Signer, error) {
	return nil, fmt.Errorf("XRPL does not support private key")
}

// DeriveFromMnemonic derives an XRPL wallet from a mnemonic and creates a LocalSigner.
func (a *Adapter) DeriveFromMnemonic(
	name, mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (wallet.Signer, error) {
	if coinType != xrplDefaultCoinType || account != 0 || index != 0 {
		return nil, fmt.Errorf("xrpl mnemonic derivation only supports m/44'/144'/0'/0/0")
	}

	mnWallet, err := xrplwallet.FromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	return NewLocalSigner(name, mnWallet), nil
}

// PersistKey stores the secret (mnemonic) in the system keyring.
func (a *Adapter) PersistKey(name string, signer wallet.Signer, secret string) error {
	kr, err := openXRPLKeyring(a.passphrase, a.homePath, a.chainName)
	if err != nil {
		return err
	}
	return setXRPLSecret(kr, a.chainName, name, secret)
}

// LoadSigner reconstructs a Signer from a persisted KeyRecord.
func (a *Adapter) LoadSigner(name string, record wallet.KeyRecord) (wallet.Signer, error) {
	switch record.Type {
	case wallet.LocalSignerType:
		kr, err := openXRPLKeyring(a.passphrase, a.homePath, a.chainName)
		if err != nil {
			return nil, err
		}
		secret, err := getXRPLSecret(kr, a.chainName, name)
		if err != nil {
			return nil, err
		}
		mnWallet, err := xrplwallet.FromMnemonic(secret)
		if err != nil {
			return nil, err
		}
		return NewLocalSigner(name, mnWallet), nil
	case wallet.RemoteSignerType:
		return NewRemoteSigner(name, record.Address, record.URL, record.Key)
	default:
		return nil, fmt.Errorf("unsupported signer type: %s for key %s", record.Type, name)
	}
}

// DeleteLocalSecret removes the secret from the system keyring.
// It is a no-op for remote signers.
func (a *Adapter) DeleteLocalSecret(name string, signer wallet.Signer) error {
	if _, ok := signer.(*LocalSigner); ok {
		kr, err := openXRPLKeyring(a.passphrase, a.homePath, a.chainName)
		if err != nil {
			return err
		}
		return deleteXRPLSecret(kr, a.chainName, name)
	}
	return nil
}
