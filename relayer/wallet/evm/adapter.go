package evm

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

const hdPathTemplate = "m/44'/%d'/%d'/0/%d"

// Adapter implements wallet.WalletAdapter for EVM/Ethereum chains.
type Adapter struct{}

var _ wallet.WalletAdapter = (*Adapter)(nil)

// LoadSigner reconstructs a Signer from a persisted KeyRecord and the secret
// retrieved from the shared keyring.
func (a *Adapter) LoadSigner(name string, record wallet.KeyRecord, secret string) (wallet.Signer, error) {
	switch record.Type {
	case wallet.PrivKeySignerType:
		privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(secret, "0x"))
		if err != nil {
			return nil, err
		}
		return NewLocalSigner(name, privateKey), nil
	case wallet.MnemonicSignerType:
		ms, err := wallet.DecodeMnemonicSecret(secret)
		if err != nil {
			return nil, err
		}
		hdw, err := hdwallet.NewFromMnemonic(ms.Mnemonic)
		if err != nil {
			return nil, err
		}
		hdPath := fmt.Sprintf(hdPathTemplate, ms.CoinType, ms.Account, ms.Index)
		derivationPath := hdwallet.MustParseDerivationPath(hdPath)
		ethAccount, err := hdw.Derive(derivationPath, true)
		if err != nil {
			return nil, err
		}
		priv, err := hdw.PrivateKey(ethAccount)
		if err != nil {
			return nil, err
		}
		return NewLocalSigner(name, priv), nil
	case wallet.RemoteSignerType:
		return NewRemoteSigner(name, record.Address, record.URL, secret)
	default:
		return nil, fmt.Errorf("unsupported signer type: %s for key %s", record.Type, name)
	}
}
