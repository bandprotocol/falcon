package xrpl

import (
	"fmt"
	"path"

	addresscodec "github.com/Peersyst/xrpl-go/address-codec"
	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"
	toml "github.com/pelletier/go-toml/v2"

	"github.com/bandprotocol/falcon/internal/os"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Wallet = &XRPLWallet{}

const (
	LocalSignerType  = "local"
	RemoteSignerType = "remote"
)

// XRPLWallet manages local and remote signers for a specific chain.
type XRPLWallet struct {
	Passphrase string
	Signers    map[string]wallet.Signer
	HomePath   string
	ChainName  string
}

// NewXRPLWallet creates a new XRPLWallet instance.
func NewXRPLWallet(passphrase, homePath, chainName string) (*XRPLWallet, error) {
	keyRecordDir := path.Join(getXRPLKeyDir(homePath, chainName)...)
	keyRecords, err := LoadKeyRecord(keyRecordDir)
	if err != nil {
		return nil, err
	}

	kr, err := openXRPLKeyring(passphrase, homePath, chainName)
	if err != nil {
		return nil, err
	}

	signers := make(map[string]wallet.Signer)
	for name, record := range keyRecords {
		var signer wallet.Signer
		switch record.Type {
		case LocalSignerType:
			secret, err := getXRPLSecret(kr, chainName, name)
			if err != nil {
				return nil, err
			}

			w, err := xrplwallet.FromSecret(secret)
			if err != nil {
				return nil, err
			}

			signer = NewLocalSigner(name, w)
		case RemoteSignerType:
			if record.Address == "" {
				return nil, fmt.Errorf("missing address for key %s", name)
			}
			if !addresscodec.IsValidClassicAddress(record.Address) {
				return nil, fmt.Errorf("invalid address: %s", record.Address)
			}

			signer = NewRemoteSigner(name, record.Address, record.Url, record.Key)
		default:
			return nil, fmt.Errorf(
				"unsupported signer type %s for chain %s, key %s",
				record.Type,
				chainName,
				name,
			)
		}

		signers[name] = signer
	}

	return &XRPLWallet{
		Passphrase: passphrase,
		Signers:    signers,
		HomePath:   homePath,
		ChainName:  chainName,
	}, nil
}

// SavePrivateKey stores the secret in keyring and writes its record.
func (w *XRPLWallet) SavePrivateKey(name string, privKey string) (addr string, err error) {
	if _, ok := w.Signers[name]; ok {
		return "", fmt.Errorf("key name exists: %s", name)
	}

	privWallet, err := xrplwallet.FromSecret(privKey)
	if err != nil {
		return
	}

	addr = privWallet.ClassicAddress.String()

	if w.IsAddressExist(addr) {
		return "", fmt.Errorf("address exists: %s", addr)
	}

	kr, err := openXRPLKeyring(w.Passphrase, w.HomePath, w.ChainName)
	if err != nil {
		return "", err
	}

	if err := setXRPLSecret(kr, w.ChainName, name, privKey); err != nil {
		return "", err
	}

	record := NewKeyRecord(LocalSignerType, "", "", nil)
	if err := w.saveKeyRecord(name, record); err != nil {
		return "", err
	}

	return addr, nil
}

// SaveRemoteSignerKey registers a remote signer under the given name.
func (w *XRPLWallet) SaveRemoteSignerKey(name, address, url string, key *string) error {
	if _, ok := w.Signers[name]; ok {
		return fmt.Errorf("key name exists: %s", name)
	}

	if !addresscodec.IsValidClassicAddress(address) {
		return fmt.Errorf("invalid address: %s", address)
	}

	if w.IsAddressExist(address) {
		return fmt.Errorf("address exists: %s", address)
	}

	record := NewKeyRecord(RemoteSignerType, address, url, key)
	if err := w.saveKeyRecord(name, record); err != nil {
		return err
	}

	return nil
}

// DeleteKey removes the signer named name, deleting its record.
func (w *XRPLWallet) DeleteKey(name string) error {
	if _, ok := w.Signers[name]; !ok {
		return fmt.Errorf("key name does not exist: %s", name)
	}

	if _, ok := w.Signers[name].(*LocalSigner); ok {
		kr, err := openXRPLKeyring(w.Passphrase, w.HomePath, w.ChainName)
		if err != nil {
			return err
		}
		if err := deleteXRPLSecret(kr, w.ChainName, name); err != nil {
			return err
		}
	}

	if err := w.deleteKeyRecord(name); err != nil {
		return err
	}

	return nil
}

// GetSigners lists all signers.
func (w *XRPLWallet) GetSigners() []wallet.Signer {
	signers := make([]wallet.Signer, 0, len(w.Signers))
	for _, signer := range w.Signers {
		signers = append(signers, signer)
	}

	return signers
}

// GetSigner returns the signer with the given name and a flag indicating if it was found.
func (w *XRPLWallet) GetSigner(name string) (wallet.Signer, bool) {
	signer, ok := w.Signers[name]
	return signer, ok
}

// IsAddressExist returns true if the given address is already added.
func (w *XRPLWallet) IsAddressExist(address string) bool {
	for _, signer := range w.Signers {
		if signer.GetAddress() == address {
			return true
		}
	}
	return false
}

// saveKeyRecord writes the KeyRecord to the file.
func (w *XRPLWallet) saveKeyRecord(name string, record KeyRecord) error {
	b, err := toml.Marshal(record)
	if err != nil {
		return err
	}

	return os.Write(b, append(getXRPLKeyDir(w.HomePath, w.ChainName), fmt.Sprintf("%s.toml", name)))
}

// deleteKeyRecord deletes the KeyRecord file.
func (w *XRPLWallet) deleteKeyRecord(name string) error {
	dir := path.Join(getXRPLKeyDir(w.HomePath, w.ChainName)...)
	filePath := path.Join(dir, fmt.Sprintf("%s.toml", name))
	return os.DeletePath(filePath)
}
