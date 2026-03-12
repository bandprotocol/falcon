package geth

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

const hdPathTemplate = "m/44'/%d'/%d'/0/%d"

// Adapter implements wallet.WalletAdapter for EVM/Ethereum chains.
type Adapter struct {
	passphrase string
	store      *keystore.KeyStore
}

var _ wallet.WalletAdapter = (*Adapter)(nil)

// DeriveFromPrivateKey parses the hex private key, derives the address and creates
// a LocalSigner without persisting the key to the keystore.
func (a *Adapter) DeriveFromPrivateKey(name, secret string) (wallet.Signer, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(secret, "0x"))
	if err != nil {
		return nil, err
	}
	return NewLocalSigner(name, privateKey), nil
}

// DeriveFromMnemonic derives an ECDSA key from a mnemonic via the BIP-44 HD path
// and creates a LocalSigner without persisting the key to the keystore.
func (a *Adapter) DeriveFromMnemonic(
	name, mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (wallet.Signer, error) {
	hdWallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	hdPath := fmt.Sprintf(hdPathTemplate, coinType, account, index)
	derivationPath := hdwallet.MustParseDerivationPath(hdPath)
	ethAccount, err := hdWallet.Derive(derivationPath, true)
	if err != nil {
		return nil, err
	}

	privHex, err := hdWallet.PrivateKeyHex(ethAccount)
	if err != nil {
		return nil, err
	}

	return a.DeriveFromPrivateKey(name, privHex)
}

// PersistKey imports the ECDSA key into the Geth keystore.
func (a *Adapter) PersistKey(name string, signer wallet.Signer, secret string) error {
	ls, ok := signer.(*LocalSigner)
	if !ok {
		return nil
	}
	_, err := a.store.ImportECDSA(ls.privateKey, a.passphrase)
	return err
}

// LoadSigner reconstructs a Signer from a persisted KeyRecord.
func (a *Adapter) LoadSigner(name string, record wallet.KeyRecord) (wallet.Signer, error) {
	switch record.Type {
	case wallet.LocalSignerType:
		gethAddr := common.HexToAddress(record.Address)
		acc := accounts.Account{Address: gethAddr}
		b, err := a.store.Export(acc, a.passphrase, a.passphrase)
		if err != nil {
			return nil, err
		}
		gethKey, err := keystore.DecryptKey(b, a.passphrase)
		if err != nil {
			return nil, err
		}
		return NewLocalSigner(name, gethKey.PrivateKey), nil
	case wallet.RemoteSignerType:
		return NewRemoteSigner(name, record.Address, record.Url, record.Key)
	default:
		return nil, fmt.Errorf("unsupported signer type: %s for key %s", record.Type, name)
	}
}

// DeleteLocalSecret removes the private key from the Geth keystore.
// It is a no-op for remote signers.
func (a *Adapter) DeleteLocalSecret(name string, signer wallet.Signer) error {
	if _, ok := signer.(*LocalSigner); ok {
		addr := common.HexToAddress(signer.GetAddress())
		return a.store.Delete(accounts.Account{Address: addr}, a.passphrase)
	}
	return nil
}
