package geth

import (
	"fmt"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

const hdPathTemplate = "m/44'/%d'/%d'/0/%d"

// GethWallet is a type alias for wallet.BaseWallet.
type GethWallet = wallet.BaseWallet

// GethAdapter implements wallet.WalletAdapter for EVM/Ethereum chains.
type GethAdapter struct {
	passphrase string
	homePath   string
	chainName  string
	store      *keystore.KeyStore
}

var _ wallet.WalletAdapter = (*GethAdapter)(nil)

// NewGethWallet creates a new GethWallet for the given chain.
func NewGethWallet(passphrase, homePath, chainName string) (*GethWallet, error) {
	keyStoreDir := path.Join(getEVMKeyStoreDir(homePath, chainName)...)
	store := keystore.NewKeyStore(keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	adapter := &GethAdapter{
		passphrase: passphrase,
		homePath:   homePath,
		chainName:  chainName,
		store:      store,
	}

	return wallet.NewBaseWallet(homePath, chainName, adapter)
}

// NormalizeAddress returns the EIP-55 checksummed form of a hex address, or an error if invalid.
func (a *GethAdapter) NormalizeAddress(addr string) (string, error) {
	if !common.IsHexAddress(addr) {
		return "", fmt.Errorf("invalid address: %s", addr)
	}
	return common.HexToAddress(addr).Hex(), nil
}

// DeriveFromPrivateKey parses the hex private key, derives the address and creates
// a LocalSigner without persisting the key to the keystore.
func (a *GethAdapter) DeriveFromPrivateKey(name, secret string) (wallet.Signer, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(secret, "0x"))
	if err != nil {
		return nil, err
	}
	return NewLocalSigner(name, privateKey), nil
}

// DeriveFromMnemonic derives an ECDSA key from a mnemonic via the BIP-44 HD path
// and creates a LocalSigner without persisting the key to the keystore.
func (a *GethAdapter) DeriveFromMnemonic(
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
func (a *GethAdapter) PersistKey(name string, signer wallet.Signer, secret string) error {
	ls, ok := signer.(*LocalSigner)
	if !ok {
		return nil
	}
	_, err := a.store.ImportECDSA(ls.privateKey, a.passphrase)
	return err
}

// LoadSigner reconstructs a Signer from a persisted KeyRecord.
func (a *GethAdapter) LoadSigner(name string, record wallet.KeyRecord) (wallet.Signer, error) {
	switch record.Type {
	case wallet.LocalSignerType:
		gethAddr, err := HexToETHAddress(record.Address)
		if err != nil {
			return nil, err
		}
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
		gethAddr, err := HexToETHAddress(record.Address)
		if err != nil {
			return nil, err
		}
		return NewRemoteSigner(name, gethAddr, record.Url, record.Key)
	default:
		return nil, fmt.Errorf("unsupported signer type: %s for key %s", record.Type, name)
	}
}

// DeleteLocalSecret removes the private key from the Geth keystore.
// It is a no-op for remote signers.
func (a *GethAdapter) DeleteLocalSecret(name string, signer wallet.Signer) error {
	if _, ok := signer.(*LocalSigner); ok {
		addr := common.HexToAddress(signer.GetAddress())
		return a.store.Delete(accounts.Account{Address: addr}, a.passphrase)
	}
	return nil
}
