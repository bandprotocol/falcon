package xrpl

import (
	"fmt"

	addresscodec "github.com/Peersyst/xrpl-go/address-codec"
	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

const (
	SaveMethodSeed     = "seed"
	SaveMethodMnemonic = "mnemonic"

	xrplDefaultCoinType = 144
)

// XRPLWallet is a type alias for wallet.BaseWallet.
type XRPLWallet = wallet.BaseWallet

// XRPLAdapter implements wallet.WalletAdapter for XRPL chains.
type XRPLAdapter struct {
	passphrase string
	homePath   string
	chainName  string
}

var _ wallet.WalletAdapter = (*XRPLAdapter)(nil)

// NewXRPLWallet creates a new XRPLWallet for the given chain.
func NewXRPLWallet(passphrase, homePath, chainName string) (*XRPLWallet, error) {
	adapter := &XRPLAdapter{
		passphrase: passphrase,
		homePath:   homePath,
		chainName:  chainName,
	}

	return wallet.NewBaseWallet(adapter)
}

// GetKeyDir returns the metadata directory path for key records.
func (a *XRPLAdapter) GetKeyDir() []string {
	return getXRPLKeyDir(a.homePath, a.chainName)
}

// ValidateAddress checks if the address is a valid XRPL classic address.
func (a *XRPLAdapter) ValidateAddress(address string) error {
	if !addresscodec.IsValidClassicAddress(address) {
		return fmt.Errorf("invalid address: %s", address)
	}
	return nil
}

// CompareAddresses returns true if two XRPL addresses are equal.
func (a *XRPLAdapter) CompareAddresses(addrA, addrB string) bool {
	return addrA == addrB
}

// DeriveFromPrivateKey is not supported for XRPL.
func (a *XRPLAdapter) DeriveFromPrivateKey(name, privateKey string) (string, wallet.Signer, error) {
	return "", nil, fmt.Errorf("XRPL does not support private key")
}

// DeriveFromMnemonic derives an XRPL wallet from a mnemonic and creates a LocalSigner.
func (a *XRPLAdapter) DeriveFromMnemonic(
	name, mnemonic string,
	coinType uint32,
	account uint,
	index uint,
) (string, wallet.Signer, string, error) {
	if coinType != xrplDefaultCoinType || account != 0 || index != 0 {
		return "", nil, "", fmt.Errorf("xrpl mnemonic derivation only supports m/44'/144'/0'/0/0")
	}

	mnWallet, err := xrplwallet.FromMnemonic(mnemonic)
	if err != nil {
		return "", nil, "", err
	}

	addr := mnWallet.ClassicAddress.String()
	return addr, NewLocalSigner(name, mnWallet), SaveMethodMnemonic, nil
}

// PersistKey stores the secret (mnemonic) in the system keyring.
func (a *XRPLAdapter) PersistKey(name string, signer wallet.Signer, secret string) error {
	kr, err := openXRPLKeyring(a.passphrase, a.homePath, a.chainName)
	if err != nil {
		return err
	}
	return setXRPLSecret(kr, a.chainName, name, secret)
}

// LoadLocalSigner reconstructs a LocalSigner from the keyring using the persisted record.
func (a *XRPLAdapter) LoadLocalSigner(name string, record wallet.KeyRecord) (wallet.Signer, error) {
	kr, err := openXRPLKeyring(a.passphrase, a.homePath, a.chainName)
	if err != nil {
		return nil, err
	}

	secret, err := getXRPLSecret(kr, a.chainName, name)
	if err != nil {
		return nil, err
	}

	var wptr *xrplwallet.Wallet
	var w xrplwallet.Wallet
	switch record.SaveMethod {
	case SaveMethodMnemonic:
		wptr, err = xrplwallet.FromMnemonic(secret)
	case SaveMethodSeed:
		w, err = xrplwallet.FromSecret(secret)
		wptr = &w
	default:
		return nil, fmt.Errorf("unsupported save method %s for key %s", record.SaveMethod, name)
	}
	if err != nil {
		return nil, err
	}

	return NewLocalSigner(name, wptr), nil
}

// NewRemoteSigner creates a new XRPL RemoteSigner with a gRPC connection to the KMS.
func (a *XRPLAdapter) NewRemoteSigner(name, address, url string, key *string) (wallet.Signer, error) {
	if address == "" {
		return nil, fmt.Errorf("missing address for key %s", name)
	}
	if !addresscodec.IsValidClassicAddress(address) {
		return nil, fmt.Errorf("invalid address: %s", address)
	}
	return NewRemoteSigner(name, address, url, key)
}

// DeleteLocalSecret removes the secret from the system keyring.
// It is a no-op for remote signers.
func (a *XRPLAdapter) DeleteLocalSecret(name string, signer wallet.Signer) error {
	if _, ok := signer.(*LocalSigner); ok {
		kr, err := openXRPLKeyring(a.passphrase, a.homePath, a.chainName)
		if err != nil {
			return err
		}
		return deleteXRPLSecret(kr, a.chainName, name)
	}
	return nil
}
