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

	return wallet.NewBaseWallet(adapter)
}

// GetKeyDir returns the metadata directory path for key records.
func (a *GethAdapter) GetKeyDir() []string {
	return getEVMMetadataDir(a.homePath, a.chainName)
}

// ValidateAddress checks if the address is a valid hex address.
func (a *GethAdapter) ValidateAddress(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("invalid address: %s", address)
	}
	return nil
}

// CompareAddresses returns true if two hex addresses refer to the same account.
func (a *GethAdapter) CompareAddresses(addrA, addrB string) bool {
	return common.HexToAddress(addrA).Cmp(common.HexToAddress(addrB)) == 0
}

// DeriveFromPrivateKey parses the hex private key, derives the address and creates
// a LocalSigner without persisting the key to the keystore.
func (a *GethAdapter) DeriveFromPrivateKey(name, secret string) (string, wallet.Signer, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(secret, "0x"))
	if err != nil {
		return "", nil, err
	}

	addr := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return addr, NewLocalSigner(name, privateKey), nil
}

// DeriveFromMnemonic derives an ECDSA key from a mnemonic via the BIP-44 HD path
// and creates a LocalSigner without persisting the key to the keystore.
func (a *GethAdapter) DeriveFromMnemonic(
	name, mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (string, wallet.Signer, string, error) {
	hdWallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", nil, "", err
	}

	hdPath := fmt.Sprintf(hdPathTemplate, coinType, account, index)
	derivationPath := hdwallet.MustParseDerivationPath(hdPath)
	ethAccount, err := hdWallet.Derive(derivationPath, true)
	if err != nil {
		return "", nil, "", err
	}

	privHex, err := hdWallet.PrivateKeyHex(ethAccount)
	if err != nil {
		return "", nil, "", err
	}

	addr, signer, err := a.DeriveFromPrivateKey(name, privHex)
	return addr, signer, "", err
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

// LoadLocalSigner reconstructs a LocalSigner by exporting the key from the keystore.
func (a *GethAdapter) LoadLocalSigner(name string, record wallet.KeyRecord) (wallet.Signer, error) {
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
}

// NewRemoteSigner creates a new RemoteSigner with a gRPC connection to the KMS.
func (a *GethAdapter) NewRemoteSigner(name, address, url string, key *string) (wallet.Signer, error) {
	gethAddr, err := HexToETHAddress(address)
	if err != nil {
		return nil, err
	}
	return NewRemoteSigner(name, gethAddr, url, key)
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
